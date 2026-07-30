package handlers

import (
	"context"

	"github.com/IvanDrf/work-hunter/pkg/common"
)

func (h *Handler) Health(ctx context.Context, _ *common.Empty) (*common.ServiceStatus, error) {
	return &common.ServiceStatus{
		Status: common.Status_AVAILABLE,
	}, nil
}
