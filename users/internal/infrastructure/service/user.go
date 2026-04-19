package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"sync"

	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	repository "github.com/IvanDrf/work-hunter/users/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/users/internal/domain/rules"
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

func (s *UserService) CreateProfile(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/CreateProfile"))

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
	log.Debug("user model created successfully", slog.String("id", user.ID.String()))

	if err := s.repo.CreateUser(ctx, user); err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user created successfully", slog.String("id", user.ID.String()))

	return modelToResp(user, log)
}

func (s *UserService) GetProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/GetProfile"))

	uuid, err := parseUUID(id, log)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))

		return nil, err
	}
	log.Info("user found successfully", slog.String("id", user.ID.String()))

	return modelToResp(user, log)
}

func (s *UserService) GetProfileByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/GetProfileByUsername"))

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Error("failed to get user by username", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user found successfully", slog.String("id", user.ID.String()))

	return modelToResp(user, log)
}

func (s *UserService) UpdateProfile(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("sscope", "infrastructure/service/UpdateProfile"))

	uuid, err := parseUUID(req.ID, log)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("user found successfully", slog.String("id", user.ID.String()))

	user.UpdateUser(req.FirstName, req.LastName, req.PhoneNumber, req.AvatarURL, req.Metadata)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user updated successfully", slog.String("id", user.ID.String()))

	return modelToResp(user, log)
}

func (s *UserService) DeleteProfile(ctx context.Context, id string) error {
	log := s.log.With(slog.String("scope", "infrastructure/service/DeleteProfile"))

	uuid, err := parseUUID(id, log)
	if err != nil {
		return err
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		log.Error("failed to get user by id", slog.String("error", err.Error()))
		return err
	}
	log.Debug("user successfully found", slog.String("id", user.ID.String()))

	if user.Status == "deleted" {
		err = s.repo.DeleteUser(ctx, uuid, true)
	} else {
		err = s.repo.DeleteUser(ctx, uuid, false)
	}

	if err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		return err
	}
	log.Info("user deleted successfully", slog.String("id", user.ID.String()))

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/ListUsers"))

	enabledFields := map[string]struct{}{
		"id": {}, "username": {}, "email": {}, "first_name": {}, "last_name": {},
		"created_at": {}, "updated_at": {},
	}

	params := make(map[string]string)
	if req.Role != "" {
		params["role"] = req.Role
	}

	if req.Status != "" {
		params["status"] = req.Status
	}

	if req.PageSize == 0 {
		params["limit"] = "100"
	} else {
		params["limit"] = strconv.Itoa(int(req.PageSize))
	}

	params["offset"] = strconv.Itoa(int(req.Offset))

	if req.SortBy == "" {
		params["order_by"] = req.SortBy
	}

	values, err := url.ParseQuery(req.SearchQuery)
	if err != nil {
		log.Error("failed to parse query", slog.String("error", err.Error()))
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to parse search query: %v", err),
			Code:    models.ErrCodeInvalidRequest,
		}
	}

	for key, val := range values {
		_, ok := enabledFields[key]
		if !ok {
			log.Error("cannot use this field", slog.String("field", key))
			return nil, &models.Error{
				Message: fmt.Sprintf("cannot use this field: %s", key),
				Code:    models.ErrCodeInvalidRequest,
			}
		}

		if len(val) > 0 {
			params[key] = val[0]
		}
	}
	log.Debug("params created successfully")

	users, totalCount, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		log.Error("failed to list users", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("users listed successfully", slog.Int("count", len(users)))

	usersResp := make([]*dto.UserResponse, 0, len(users))
	for _, val := range users {
		userResp, err := modelToResp(val, log)
		if err != nil {
			return nil, err
		}

		usersResp = append(usersResp, userResp)
	}

	var hasNextPage bool
	if len(users) < int(totalCount) {
		hasNextPage = true
	}

	return &dto.ListUsersResponse{
		Users:      usersResp,
		TotalCount: totalCount,
		HasNext:    hasNextPage,
	}, nil
}

// TODO: refactor
func (s *UserService) UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/UpdateUserStatus"))

	id, err := parseUUID(req.ID, log)
	if err != nil {
		return nil, err
	}

	user := new(models.User)
	var wg sync.WaitGroup
	wg.Go(func() {
		user, errGet := s.repo.GetUserByID(ctx, id)
		if err != nil {
			log.Error("failed to get user", slog.String("error", errGet.Error()))
		}
		log.Debug("user found successfully", slog.String("id", user.ID.String()))

		user.Status = rules.UserStatus(req.Status)
	})

	if err = s.repo.UpdateUserStatus(ctx, id, rules.UserStatus(req.Status)); err != nil {
		log.Error("failed to update user status", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("user status updated successfully")

	wg.Wait()

	return modelToResp(user, log)
}

func (s *UserService) Close() {
	s.repo.Close()
	s.log.Info("connection closed successfully")
}

func parseUUID(id string, log *slog.Logger) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse uuid from string", slog.String("error", err.Error()))
		return uuid, &models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInvalidRequest,
		}
	}
	log.Debug("uuid parsed successfully")
	return uuid, nil
}

func modelToResp(user *models.User, log *slog.Logger) (*dto.UserResponse, error) {
	metadata := make(map[string]string)
	if err := json.Unmarshal(user.Metadata, &metadata); err != nil {
		log.Error("failed to unmarshal json data", slog.String("error", err.Error()))
		return nil, &models.Error{
			Message: fmt.Sprintf("failed to unmarshal json data: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}
	log.Debug("user converted successfully", slog.String("id", user.ID.String()))

	return &dto.UserResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
		Status:      string(user.Status),
		Role:        string(user.Role),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Metadata:    metadata,
	}, nil
}
