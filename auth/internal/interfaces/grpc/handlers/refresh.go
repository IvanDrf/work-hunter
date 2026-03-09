package handlers

import (
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) RefreshTokens(tokens *auth_api.JwtTokens) (*auth_api.JwtTokens, error) {
	access, refresh, err := h.authService.RefreshTokens(tokens.Access, tokens.Refresh)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrInvalidJWT:
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
