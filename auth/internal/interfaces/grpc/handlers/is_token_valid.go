package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) IsTokenValid(ctx context.Context, access *auth_api.AccessToken) (*auth_api.TokenPayload, error) {
	slog.Info("IsTokenValid got request")

	payload, err := h.authService.GetTokenPayload(ctx, access.GetAccess())
	if err != nil {
		return nil, handleError(err, "IsTokenValid error")
	}

	slog.Info("IsTokenValid successful response")
	return &auth_api.TokenPayload{
		Id:          payload.UserID,
		Verificated: payload.Verificated,
		Role:        auth_api.Role(auth_api.Role_value[string(payload.Role)]),
	}, nil
}
