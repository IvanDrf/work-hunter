package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Register(ctx context.Context, user *auth_api.User) (*auth_api.JwtTokens, error) {
	access, refresh, err := h.authService.RegisterUser(ctx, user.Email, user.Password)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, e.Message)

		case models.ErrCodeInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	err = h.verificationService.SendVerificationEmail(ctx, user.Email)
	if err != nil {
		// TODO: add logging
	}

	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
