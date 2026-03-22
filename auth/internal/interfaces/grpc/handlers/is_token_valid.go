package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) IsTokenValid(ctx context.Context, access *auth_api.AccessToken) (*auth_api.TokenPayload, error) {
	payload, err := h.authService.GetTokenPayload(ctx, access.Access)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInvalidJWT:
			return nil, status.Error(codes.InvalidArgument, e.Message)
		}
	}

	return &auth_api.TokenPayload{
		Id:          payload.ID.String(),
		Verificated: payload.Verificated,
	}, nil
}
