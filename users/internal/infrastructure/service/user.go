package service

import (
	"context"
	"fmt"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	repository "github.com/IvanDrf/work-hunter/users/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
	"github.com/google/uuid"
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
	uuid, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	user := models.NewUser(uuid, req.Username, req.Email, req.FirstName, req.LastName, req.PhoneNumber)
	_, err = s.repo.GetUserByID(ctx, uuid)
	if err == nil {
		return nil, &models.Error{
			Message: "user already exist",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*models.User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error) {
	uuid, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	user.UpdateUser(req.FirstName, req.LastName, req.PhoneNumber, req.AvatarURL, req.Metadata)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteProfile(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return err
	}

	if user.Status == "deleted" {
		err = s.repo.DeleteUser(ctx, uuid, true)
	} else {
		err = s.repo.DeleteUser(ctx, uuid, false)
	}

	if err != nil {
		return err
	}
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
