package repo

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUser(ctx context.Context, email string) (*models.User, error)
}
