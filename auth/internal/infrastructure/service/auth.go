package service

import (
	"context"

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

func (a *AuthService) RegisterUser(ctx context.Context, email string, password string) (string, string, error) {
	_, err := a.userRepo.FindUser(ctx, email)
	if err == nil {
		return "", "", models.Error{
			Message: "user with that email already exists",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	user, err := models.NewUser(email, password)
	if err != nil {
		return "", "", err
	}

	err = a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", "", models.Error{
			Message: "can't register new user",
			Code:    models.ErrCodeInternal,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(user.ID, user.Verificated)
	if err != nil {
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	return access, refresh, nil
}

func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (string, string, error) {
	user, err := a.userRepo.FindUser(ctx, email)
	if err != nil {
		return "", "", models.Error{
			Message: "user with that email doesn't exists",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !rules.IsPasswordsSame(password, user.HashedPassword) {
		return "", "", models.Error{
			Message: "invalid password",
			Code:    models.ErrCodeInvalidPassword,
		}
	}

	access, refresh, err := a.jwter.CreateTokens(user.ID, user.Verificated)
	if err != nil {
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	return access, refresh, nil
}

func (a *AuthService) RefreshTokens(ctx context.Context, refresh string) (string, string, error) {
	access, refresh, err := a.jwter.RefreshTokens(refresh)
	if err != nil {
		return "", "", models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	return access, refresh, nil
}
