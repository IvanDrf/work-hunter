package repoPort

import (
	"context"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, updatedUser *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID, permanent bool) error
	ListUsers(ctx context.Context, params map[string]string) ([]*models.User, int32, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status rules.UserStatus) error
}
