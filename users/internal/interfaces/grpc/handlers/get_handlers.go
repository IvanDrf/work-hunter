package handlers

import (
	"context"
	"errors"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func convertUserResponseToUserProfile(user *dto.UserResponse) *user_api.UserProfile {
	return &user_api.UserProfile{
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Avatar_URL:  user.AvatarURL,
		Status:      user_api.UserStatus(user_api.UserStatus_value[user.Status]),
		Role:        user_api.UserRole(user_api.UserRole_value[user.Role]),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
		Metadata:    user.Metadata,
	}
}
