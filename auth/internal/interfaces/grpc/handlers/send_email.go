package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) SendVerificationEmail(ctx context.Context, token *auth_api.AccessToken) (*auth_api.Empty, error) {
	slog.Info("SendVerificationEmail got request")

	err := h.verificationService.ResendVerificationEmail(ctx, token.GetAccess())
	if err != nil {
		return nil, handleError(err, "SendVerificationEmail error")
	}

	slog.Info("SendVerificationEmail successful response")
	return &auth_api.Empty{}, nil
}
