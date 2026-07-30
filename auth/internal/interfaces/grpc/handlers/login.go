package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) Login(ctx context.Context, user *auth_api.User) (*auth_api.JwtTokens, error) {
	slog.Info("Login got request")

	access, refresh, err := h.authService.LoginUser(ctx, user.GetEmail(), user.GetPassword())
	if err != nil {
		return nil, handleError(err, "Login error")
	}

	slog.Info("Login successful response")
	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
