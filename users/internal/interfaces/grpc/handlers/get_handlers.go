package handlers

import (
	"context"
	"errors"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetProfile(ctx context.Context, req *user_api.GetProfileRequest) (*user_api.UserProfile, error) {
	log := h.log.With(slog.String("scope", "interfaces/grpc/handlers/GetProfile"))

	log.Info("GetProfile got request")
	user, err := h.UserService.GetProfile(ctx, req.UserId)

	var e models.Error
	if errors.As(err, &e) {
		log.Error("Get profile error", slog.String("error", e.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)
		default:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	log.Info("Get profile successfully response")
	return convertUserResponseToUserProfile(user), nil
}

func (h *Handler) GetProfileByUsername(ctx context.Context, req *user_api.GetProfileByUsernameRequest) (*user_api.UserProfile, error) {
	log := h.log.With(slog.String("scope", "interfaces/grpc/handlers/GetProfileByUsername"))

	log.Info("GetProfileByUsername got request")
	user, err := h.UserService.GetProfileByUsername(ctx, req.Username)

	var e models.Error
	if errors.As(err, &e) {
		log.Error("Get profile by username error", slog.String("error", e.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)
		default:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	log.Info("Get profile by username successfully response")
	return convertUserResponseToUserProfile(user), nil
}
