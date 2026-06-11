package handlers

import (
	"context"
	"errors"
	"log/slog"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/domain/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DeleteProfile(ctx context.Context, req *user_api.DeleteProfileRequest) (*emptypb.Empty, error) {
	log := h.log.With(slog.String("scope", "interfaces/grpc/handlers/DeleteProfile"))
	log.Info("DeleteProfile got request")

	err := h.UserService.DeleteProfile(ctx, req.UserId)
	var e models.Error
	if errors.As(err, &e) {
		log.Error("Delete profile error", slog.String("error", e.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)
		default:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	log.Info("Delete profile successfully response")
	return nil, nil
}
