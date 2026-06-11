package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ChangePassword(ctx context.Context, req *auth_api.ChangePasswordRequest) (*auth_api.Empty, error) {
	slog.Info("ChangePassword got request")

	err := h.authService.ChangeUserPassword(ctx, req.Access, req.Old, req.New)

	var e models.Error
	if errors.As(err, &e) {
		slog.Error("ChangePassword error", slog.String("error", err.Error()))

		switch e.Code {
		case models.ErrCodeInvalidJWT, models.ErrCodeInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	slog.Info("ChangePassword successfull response")
	return &auth_api.Empty{}, nil
}
