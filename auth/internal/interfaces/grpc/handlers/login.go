package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Login(ctx context.Context, user *auth_api.User) (*auth_api.JwtTokens, error) {
	access, refresh, err := h.authService.LoginUser(ctx, user.Username, user.Password)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrCodeInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
