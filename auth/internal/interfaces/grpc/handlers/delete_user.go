package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) DeleteUser(ctx context.Context, req *auth_api.DeleteUserRequest) (*auth_api.Empty, error) {
	slog.Info("DeleteUser got request")

	err := h.authService.DeleteUser(ctx, req.GetAccess(), req.GetPassword())
	if err != nil {
		return nil, handleError(err, "DeleteUser error")
	}

	slog.Info("DeleteUser successful response")
	return &auth_api.Empty{}, nil
}
