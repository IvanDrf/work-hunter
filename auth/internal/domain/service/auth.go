package service

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, username string, password string) (string, string, error)
	LoginUser(ctx context.Context, username string, password string) (string, string, error)

	RefreshTokens(ctx context.Context, access string, refresh string) (string, string, error)
}
