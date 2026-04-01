package service

import (
	"context"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	repository "github.com/IvanDrf/work-hunter/users/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
)

type UserService struct {
	repo repository.UserRepository
	log  *logger.Logger
}

func NewUserService(repo repository.UserRepository, log *logger.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) CreateProfile(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	// TODO
	return nil, nil
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*models.User, error) {
	//TODO
	return nil, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	//TODO
	return nil, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error) {
	// TODO
	return nil, nil
}

func (s *UserService) DeleteProfile(ctx context.Context, id string) error {
	// TODO
	return nil
}

func (s *UserService) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	// TODO
	return nil, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error) {
	// TODO
	return nil, nil
}
