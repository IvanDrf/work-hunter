package ports

import (
	"context"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
)

type AuthClient interface {
	SendRegisterRequest(ctx context.Context, email string, password string, role models.UserRole) (*models.Tokens, error)

	Close()
}
