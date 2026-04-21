package handlers

import (
	"context"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
)

func (h *Handler) Health(ctx context.Context, _ *auth_api.Empty) (*auth_api.ServiceStatus, error) {
	return &auth_api.ServiceStatus{
		Code: int32(codes.OK),
	}, nil
}
