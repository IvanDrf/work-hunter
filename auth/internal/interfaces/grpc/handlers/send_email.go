package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) SendVerificationEmail(ctx context.Context, email *auth_api.Email) (*auth_api.AcceptStatus, error) {
	slog.Info("SendVerificationEmail got request")
	err := h.verificationService.SendVerificationEmail(ctx, email.Email)

	var e models.Error
	if errors.As(err, &e) {
		slog.Error("SendVerificationEmail error", slog.String("error", err.Error()))

		switch e.Code {
		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	slog.Info("SendVerificationEmail successfull response")
	return &auth_api.AcceptStatus{
		Accepted: true,
	}, nil
}
