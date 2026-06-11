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

func (h *Handler) CreateProfile(ctx context.Context, req *user_api.CreateProfileRequest) (*user_api.UserProfile, error) {
	log := h.log.With(slog.String("scope", "interfaces/grpc/handlers/CreateProfile"))

	log.Info("CreateProfile got request")

	user, err := h.UserService.CreateProfile(ctx, convertCreateProfileResponseToDto(req))

	var e models.Error
	if errors.As(err, &e) {
		log.Error("Create profile error", slog.String("error", e.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		case models.ErrCodeUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, e.Message)
		default:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	log.Info("Get profile successfully response")
	return convertUserResponseToUserProfile(user), nil
}
