package handlers

import (
	"context"
	"errors"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SendVerificationEmail(ctx context.Context, email *auth_api.Email) (*auth_api.AcceptStatus, error) {
	err := h.verificationService.SendVerificationEmail(ctx, email.Email)

	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	return &auth_api.AcceptStatus{
		Accepted: true,
	}, nil
}
