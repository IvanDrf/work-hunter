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

func (h *Handler) UpdateProfile(ctx context.Context, req *user_api.UpdateProfileRequest) (*user_api.UserProfile, error) {
	log := h.log.With(slog.String("scope", "interfaces/grpc/handlers/UpdateProfile"))
	log.Info("UpdateProfile got request")

	dto, err := convertUpdateProfileResponseToDto(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err := h.UserService.UpdateProfile(ctx, dto)

	var e models.Error
	if errors.As(err, &e) {
		log.Error("Update profile error", slog.String("error", e.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)
		default:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	log.Info("Update profile successfully response")
	return convertUserResponseToUserProfile(user), nil
}
