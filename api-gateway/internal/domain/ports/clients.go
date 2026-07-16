package ports

import (
	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
)

type AuthClient interface {
	Health(ctx context.Context)

	SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (*models.Tokens, error)
	SendLoginRequest(ctx context.Context, email string, password string) (*models.Tokens, error)
	SendChangePasswordRequest(ctx context.Context, access string, oldPassword string, newPassword string) error
	SendDeleteUserRequest(ctx context.Context, access string, password string) error

	SendVerificationEmailRequest(ctx context.Context, access string) error
	SendVerifyEmailRequest(ctx context.Context, token string) (*models.Tokens, error)

	SendRefreshTokensRequest(ctx context.Context, refresh string) (*models.Tokens, error)
	SendIsTokenValidRequest(ctx context.Context, access string) (*models.TokenPayload, error)

	Close()
}
