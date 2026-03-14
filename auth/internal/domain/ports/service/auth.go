package service

import (
	"context"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email string, password string) (string, string, error)
	LoginUser(ctx context.Context, email string, password string) (string, string, error)

	RefreshTokens(ctx context.Context, refresh string) (string, string, error)
}
