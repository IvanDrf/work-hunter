package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) RefreshTokens(ctx context.Context, token *auth_api.RefreshToken) (*auth_api.JwtTokens, error) {
	slog.Info("RefreshTokens got request")

	access, refresh, err := h.authService.RefreshTokens(ctx, token.GetRefresh())
	if err != nil {
		return nil, handleError(err, "RefreshTokens error")
	}

	slog.Info("RefreshTokens successful response")
	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
