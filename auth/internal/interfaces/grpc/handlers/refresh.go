package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) RefreshTokens(ctx context.Context, tokens *auth_api.RefreshToken) (*auth_api.JwtTokens, error) {
	access, refresh, err := h.authService.RefreshTokens(ctx, tokens.Refresh)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInvalidJWT:
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
