package servicePort

import (
	"context"

	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
)

type UserService interface {
	CreateProfile(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetProfile(ctx context.Context, id string) (*dto.UserResponse, error)
	GetProfileByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteProfile(ctx context.Context, id string) error
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*dto.UserResponse, error)
}
