package repository

import (
	"context"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/models"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) ([]*models.User, int32, error)
	UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error)
}
