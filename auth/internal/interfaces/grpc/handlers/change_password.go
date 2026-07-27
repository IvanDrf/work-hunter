package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) ChangePassword(ctx context.Context, req *auth_api.ChangePasswordRequest) (*auth_api.Empty, error) {
	slog.Info("ChangePassword got request")

	err := h.authService.ChangeUserPassword(ctx, req.GetAccess(), req.GetOld(), req.GetNew())
	if err != nil {
		return nil, handleError(err, "ChangePassword error")
	}

	slog.Info("ChangePassword successful response")
	return &auth_api.Empty{}, nil
}
