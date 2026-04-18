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

func (h *Handler) SendVerificationEmail(ctx context.Context, token *auth_api.AccessToken) (*auth_api.Empty, error) {
	slog.Info("SendVerificationEmail got request")
	err := h.verificationService.ResendVerificationEmail(ctx, token.Access)

	var e models.Error
	if errors.As(err, &e) {
		slog.Error("SendVerificationEmail error", slog.String("error", err.Error()))

		switch e.Code {
		case models.ErrCodeInvalidJWT:
			return nil, status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeUserNotFound:
			return nil, status.Error(codes.NotFound, e.Message)

		case models.ErrCodeUserAlreadyVerificated:
			return nil, status.Error(codes.AlreadyExists, e.Message)

		case models.ErrCodeInternal:
			return nil, status.Error(codes.Internal, e.Message)
		}
	}

	slog.Info("SendVerificationEmail successfull response")
	return &auth_api.Empty{}, nil
}
