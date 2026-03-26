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

func (a *AuthService) RegisterUser(ctx context.Context, email string, password string) (string, string, error) {
	_, err := a.userRepo.FindUser(ctx, email)
	if err == nil {
		slog.Info("RegisterUser user with that email already exists")
		return "", "", models.Error{
			Message: "user with that email already exists",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	user, err := models.NewUser(email, password)
	if err != nil {
		slog.Error("RegisterUser error", slog.String("error", err.Error()))
		return "", "", err
	}

	err = a.userRepo.CreateUser(ctx, user)
	if err != nil {
		slog.Error("RegisterUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't register new user",
			Code:    models.ErrCodeInternal,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(user.ID, user.Verificated)
	if err != nil {
		slog.Error("RegisterUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("RegisterUser success")
	return access, refresh, nil
}

func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (string, string, error) {
	user, err := a.userRepo.FindUser(ctx, email)
	if err != nil {
		slog.Error("LoginUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "user with that email doesn't exists",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !rules.IsPasswordsSame(password, user.HashedPassword) {
		slog.Info("LoginUser incorrect password")
		return "", "", models.Error{
			Message: "invalid password",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(user.ID, user.Verificated)
	if err != nil {
		slog.Error("LoginUser error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("LoginUser success")
	return access, refresh, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, refresh string) (string, string, error) {
	access, refresh, err := a.jwter.RefreshTokens(refresh)
	if err != nil {
		slog.Error("RefreshTokens error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	slog.Info("RefreshTokens success")
	return access, refresh, nil
}

func (a *AuthService) GetTokenPayload(ctx context.Context, access string) (*models.JwtPayload, error) {
	id, verificated, err := a.jwter.GetPayload(access)
	if err != nil {
		slog.Error("GetTokenPayload error", slog.String("error", err.Error()))
		return nil, models.Error{
			Code:    models.ErrCodeInvalidJWT,
			Message: "invalid jwt token",
		}
	}

	slog.Info("GetTokenPayload success")
	return &models.JwtPayload{
		ID:          id,
		Verificated: verificated,
	}, nil
}
