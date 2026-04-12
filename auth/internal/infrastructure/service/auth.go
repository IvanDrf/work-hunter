package service

import (
	"context"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/google/uuid"
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

func (a *AuthService) DeleteUser(ctx context.Context, access string, password string) error {
	payload, err := a.jwter.GetPayload(access)
	if err != nil {
		slog.Info("auth:DeleteUser got invalid jwt token")
		return models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		slog.Info("auth:DeleteUser got invalid userID in jwt token")
		return models.Error{
			Message: "invalid userID in jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	user, err := a.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		slog.Info("auth:DeleteUser can't find user in database with given userID")
		return models.Error{
			Message: "can't find user with given userID in jwt token",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !rules.IsPasswordsSame(password, user.HashedPassword) {
		slog.Info("auth:DeleteUser password in request and password in database are different, can't change password")
		return models.Error{
			Message: "password is incorrect",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	err = a.userRepo.DeleteUser(ctx, user.Email)
	if err != nil {
		slog.Error("auth:DeleteUser can't delete user from databse", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't delete user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("auth:DeleteUser success")
	return nil
}

func (a *AuthService) ChangeUserPassword(ctx context.Context, access string, old string, new string) error {
	payload, err := a.jwter.GetPayload(access)
	if err != nil {
		slog.Info("auth:ChangeUserPassword got invalid jwt token")
		return models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		slog.Info("auth:ChangeUserPassword got invalid userID in jwt token")
		return models.Error{
			Message: "invalid userID in jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	user, err := a.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		slog.Info("auth:ChangeUserPassword can't find user with given userID in jwt token")
		return models.Error{
			Message: "can't find user with given userID",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !rules.IsPasswordsSame(old, user.HashedPassword) {
		slog.Info("auth:ChangeUserPassword old password in request and password in dataabse are different")
		return models.Error{
			Message: "password is incorrect",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	if !rules.IsPasswordCorrect(new) {
		slog.Info("auth:ChangeUserPassword new password is incorrect, doesn't fit rules")
		return models.Error{
			Message: "new password is incorrect",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	hashed, err := rules.HashPassword(new)
	if err != nil {
		slog.Error("auth:ChangeUserPassword can't hash new password", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't hash new password",
			Code:    models.ErrCodeInternal,
		}
	}

	err = a.userRepo.ChangeUserPassword(ctx, userID, hashed)
	if err != nil {
		slog.Error("auth:ChangeUserPassword can't change user password in database", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't change user password in database",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("auth:ChangeUserPassword success")
	return nil
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
