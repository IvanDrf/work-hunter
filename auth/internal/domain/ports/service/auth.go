package service

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string, role string) (string, string, error)
	LoginUser(ctx context.Context, email string, password string) (string, string, error)
	DeleteUser(ctx context.Context, access string, password string) error

	RefreshTokens(ctx context.Context, refresh string) (string, string, error)
	GetTokenPayload(ctx context.Context, access string) (*models.JwtPayload, error)

	Close()
}
