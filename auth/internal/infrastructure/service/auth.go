package service

import (
	"context"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
)

type AuthService struct {
	userRepo repo.UserRepo
	jwter    jwt.Jwter
}

func NewAuthService(userRepo repo.UserRepo, jwter jwt.Jwter) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwter:    jwter,
	}
}

func (a *AuthService) Close() {
	a.userRepo.Close()
}

func (a *AuthService) RegisterUser(ctx context.Context, email string, password string, role string) (string, string, error) {
	_, err := a.userRepo.FindUserByEmail(ctx, email)
	if err == nil {
		slog.Info("auth:RegisterUser user with that email already exists")
		return "", "", models.Error{
			Message: "user with that email already exists",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	userRole := models.ROLES[role]
	if userRole == "" {
		slog.Info("auth:RegisterUser invalid role", slog.String("role", role))
		return "", "", models.Error{
			Message: "invalid role",
			Code:    models.ErrCodeInvalidUserRole,
		}
	}

	user, err := models.NewUser(email, password, userRole)
	if err != nil {
		slog.Error("auth:RegisterUser error", slog.String("error", err.Error()))
		return "", "", err
	}

	err = a.userRepo.CreateUser(ctx, user)
	if err != nil {
		slog.Error("auth:RegisterUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't register new user",
			Code:    models.ErrCodeInternal,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(&models.JwtPayload{
		UserID:      user.ID.String(),
		Verificated: user.Verificated,
		Role:        user.Role,
	})

	if err != nil {
		slog.Error("auth:RegisterUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("auth:RegisterUser success")
	return access, refresh, nil
}

func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (string, string, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		slog.Error("auth:LoginUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "user with that email doesn't exists",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !rules.IsPasswordsSame(password, user.HashedPassword) {
		slog.Info("auth:LoginUser incorrect password")
		return "", "", models.Error{
			Message: "invalid password",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(&models.JwtPayload{
		UserID:      user.ID.String(),
		Verificated: user.Verificated,
		Role:        user.Role,
	})
	if err != nil {
		slog.Error("auth:LoginUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("auth:LoginUser success")
	return access, refresh, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, refresh string) (string, string, error) {
	payload, err := a.jwter.GetPayload(refresh)
	if err != nil {
		slog.Error("auth:RefreshTokens error cant get valid payload", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	if !payload.IsPayloadValid() {
		slog.Error("auth:RefreshTokens error payload isn't valid")
		return "", "", models.Error{
			Message: "invalid payload in jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(payload)
	if err != nil {
		slog.Error("auth:RefreshTokens error, cant create new tokens", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create new jwt tokens",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("auth:RefreshTokens success")
	return access, refresh, nil
}

func (a *AuthService) GetTokenPayload(ctx context.Context, access string) (*models.JwtPayload, error) {
	payload, err := a.jwter.GetPayload(access)
	if err != nil {
		slog.Error("auth:GetTokenPayload error", slog.String("error", err.Error()))
		return nil, models.Error{
			Code:    models.ErrCodeInvalidJWT,
			Message: "invalid jwt token",
		}
	}

	if !payload.IsPayloadValid() {
		slog.Error("auth:RefreshTokens error payload isn't valid")
		return nil, models.Error{
			Message: "invalid payload in jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	slog.Info("auth:GetTokenPayload success")
	return payload, nil
}
