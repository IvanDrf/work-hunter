package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) VerifyEmail(ctx context.Context, token *auth_api.VerifToken) (*auth_api.JwtTokens, error) {
	slog.Info("VerifyEmail got request")

	access, refresh, err := h.verificationService.VerifyEmailByToken(ctx, token.GetToken())
	if err != nil {
		return nil, handleError(err, "VerifyEmail error")
	}

	slog.Info("VerifyEmail successful response")
	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
