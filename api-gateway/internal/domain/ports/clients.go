package ports

import (
	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
)

type AuthClient interface {
	Health(ctx context.Context)

	SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (*models.Tokens, error)
	SendLoginRequest(ctx context.Context, email string, password string) (*models.Tokens, error)
	SendChangePasswordRequest(ctx context.Context, access string, old string, new string) error

	Close()
}
