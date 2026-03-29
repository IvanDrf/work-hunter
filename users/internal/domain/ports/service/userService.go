package service

import (
	"context"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
)

type UserService interface {
	CreateProfile(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error)
	GetProfile(ctx context.Context, id string) (*models.User, error)
	GetProfileByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error)
	DeleteProfile(ctx context.Context, id string) error
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error)
}
