package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) VerifyEmail(ctx context.Context, token *auth_api.VerifToken) (*auth_api.JwtTokens, error) {
	access, refresh, err := h.verificationService.VerifyEmailByToken(ctx, token.Token)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrOutdatedToken:
			return nil, status.Error(codes.DeadlineExceeded, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
