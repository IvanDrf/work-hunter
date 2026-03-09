package service

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/repo"
)

type AuthService struct {
	repo repo.UserRepo
}

func (a *AuthService) RegisterUser(ctx context.Context, username string, password string) (string, string, error)
