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

func (h *Handler) Login(ctx context.Context, user *auth_api.User) (*auth_api.JwtTokens, error) {
	slog.Info("Login got request")
	access, refresh, err := h.authService.LoginUser(ctx, user.Email, user.Password)

	var e models.Error
	if errors.As(err, &e) {
		slog.Error("Login error", slog.String("error", err.Error()))

		switch e.Code {
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrCodeInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	slog.Info("Login successfull response")
	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
