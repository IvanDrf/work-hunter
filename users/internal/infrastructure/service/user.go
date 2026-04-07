package service

import (
	"context"
	"fmt"
	"log/slog"

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
	log := s.log.With("scope", "infrastructure/service/CreateProfile")

	uuid, err := parseUUID(req.ID, log)
	if err != nil {
		return nil, err
	}

	user := models.NewUser(uuid, req.Username, req.Email, req.FirstName, req.LastName, req.PhoneNumber)
	_, err = s.repo.GetUserByID(ctx, uuid)
	if err == nil {
		log.Error("user already exist")

		return nil, &models.Error{
			Message: "user already exist",
			Code:    models.ErrCodeUserAlreadyExists,
		}
	}
	log.Debug("user model created successfully", "user", user)

	if err := s.repo.CreateUser(ctx, user); err != nil {
		log.Error("failed to create user", "error", err)
		return nil, err
	}
	log.Info("user created successfully", "user", user)

	return user, nil
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*models.User, error) {
	log := s.log.With("scope", "infrastructure/service/GetProfile")

	uuid, err := parseUUID(id, log)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", "error", err)

		return nil, err
	}
	log.Info("user found successfully", "user", user)

	return user, nil
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (*models.User, error) {
	log := s.log.With("scope", "infrastructure/service/GetProfileByUsername")

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("failed to get user by username", "error", err)
		return nil, err
	}
	log.Info("user found successfully", "user", user)

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*models.User, error) {
	log := s.log.With("scope", "infrastructure/service/UpdateProfile")

	uuid, err := parseUUID(req.ID, log)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", "error", err)
		return nil, err
	}
	log.Debug("user found successfully", "user", user)

	user.UpdateUser(req.FirstName, req.LastName, req.PhoneNumber, req.AvatarURL, req.Metadata)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		log.Error("failed to update user", "error", err)
		return nil, err
	}
	log.Info("user updated successfully", "user", user)

	return user, nil
}

func (s *UserService) DeleteProfile(ctx context.Context, id string) error {
	log := s.log.With("scope", "infrastructure/service/DeleteProfile")

	uuid, err := parseUUID(id, log)
	if err != nil {
		return err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", "error", err)
		return err
	}
	log.Debug("user successfully found", "user", user)

	if user.Status == "deleted" {
		err = s.repo.DeleteUser(ctx, uuid, true)
	} else {
		err = s.repo.DeleteUser(ctx, uuid, false)
	}

	if err != nil {
		log.Error("failed to delete user", "error", err)
		return err
	}
	log.Info("user deleted successfully", "user", user)

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {

	return nil, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*models.User, error) {
	// TODO
	return nil, nil
}

func parseUUID(id string, log *slog.Logger) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse uuid from string", "error", err)
		return uuid, &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}
	log.Debug("uuid parsed successfully", "uuid", uuid)
	return uuid, nil
}
