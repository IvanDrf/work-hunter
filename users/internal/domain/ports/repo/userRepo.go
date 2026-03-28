package repository

import (
	"context"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/models"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/interfaces/grpc/dto"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error)
}
