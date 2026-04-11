package repo

import (
	"context"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error

	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)

	DeleteUser(ctx context.Context, email string) error

	VerifyEmail(ctx context.Context, email string) error

	Close()
}
