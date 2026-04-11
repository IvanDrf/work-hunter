package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) DeleteUser(ctx context.Context, req *auth_api.DeleteUserRequest) (*auth_api.DeleteUserStatus, error) {
	err := h.authService.DeleteUser(ctx, req.Access, req.Password)
	var e models.Error

	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInvalidJWT, models.ErrCodeInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	return &auth_api.DeleteUserStatus{
		Deleted: true,
	}, nil
}
