package service

import (
	"context"
	"fmt"
	"log/slog"

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

	user := models.NewUser(uuid, req.Email, req.FirstName, req.LastName, req.CompanyName, req.Verificated)
	log.Debug("user model created successfully", slog.String("id", user.ID.String()))

	if err := s.repo.CreateUser(ctx, user); err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user created successfully", slog.String("id", user.ID.String()))

	return modelToResp(user), nil
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

	return modelToResp(user), nil
}

func (s *UserService) GetProfileByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/GetProfileByUsername"))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get user by username", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user found successfully", slog.String("id", user.ID.String()))

	return modelToResp(user), nil
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

	user.UpdateUser(req.FirstName, req.LastName, req.CompanyName, req.Verificated)
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user updated successfully", slog.String("id", user.ID.String()))

	return modelToResp(user), nil
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
	log := s.log.With(slog.String("scope", "service/ListUsers"))

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 100
	}

	repoParams := models.ListUsersParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}

	if req.Filter != nil {
		repoParams.Status = req.Filter.Status
		repoParams.Role = req.Filter.Role
		repoParams.SearchQuery = req.Filter.SearchQuery
	}

	users, totalCount, err := s.repo.ListUsers(ctx, &repoParams)
	if err != nil {
		log.Error("failed to list users", slog.String("error", err.Error()))
		return nil, err
	}

	usersResp := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		usersResp = append(usersResp, modelToResp(user))
	}

	totalPages := int32((totalCount + int32(req.PageSize) - 1) / int32(req.PageSize))

	log.Info("users listed successfully",
		slog.Int("count", len(users)),
		slog.Int("total", int(totalCount)),
	)

	return &dto.ListUsersResponse{
		Users:      usersResp,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, req *dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	log := s.log.With(slog.String("scope", "infrastructure/service/UpdateUserStatus"))

	id, err := parseUUID(req.ID, log)
	if err != nil {
		return nil, err
	}

	if err = s.repo.UpdateUserStatus(ctx, id, rules.UserStatus(req.Status)); err != nil {
		log.Error("failed to update user status", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("user status updated successfully")

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
	}
	log.Debug("user found successfully", slog.String("id", user.ID.String()))

	return modelToResp(user), nil
}

func (s *UserService) Close() {
	s.repo.Close()
	s.log.Info("connection closed successfully")
}

func parseUUID(id string, log *slog.Logger) (uuid.UUID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse uuid from string", slog.String("error", err.Error()))
		return uuid, models.Error{
			Message: fmt.Sprintf("failed to parse uuid from string: %v", err),
			Code:    models.ErrCodeInvalidRequest,
		}
	}
	log.Debug("uuid parsed successfully")
	return uuid, nil
}

func modelToResp(user *models.User) *dto.UserResponse {
	resp := &dto.UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		Status:      string(user.Status),
		Role:        string(user.Role),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Verificated: user.Verificated,
	}

	if user.FirstName != "" {
		resp.FirstName = user.FirstName
	}
	if user.LastName != "" {
		resp.LastName = user.LastName
	}
	if user.CompanyName != "" {
		resp.CompanyName = user.CompanyName
	}

	return resp
}
