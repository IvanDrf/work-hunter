package handlers

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error, message string) error {
	var e models.Error
	slog.Error(message, slog.String("error", err.Error()))
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeInvalidJWT, models.ErrCodeInvalidPassword, models.ErrCodeInvalidUserRole, models.ErrCodeInvalidEmail:
			return status.Error(codes.InvalidArgument, e.Message)

		case models.ErrCodeUserNotFound:
			return status.Error(codes.NotFound, e.Message)

		case models.ErrCodeUserAlreadyExists:
			return status.Error(codes.AlreadyExists, e.Message)

		case models.ErrCodeInternal:
			return status.Error(codes.Internal, e.Message)

		case models.ErrCodeUserAlreadyVerificated:
			return status.Error(codes.AlreadyExists, e.Message)

		case models.ErrCodeOutdatedToken:
			return status.Error(codes.DeadlineExceeded, e.Message)
		}
	}

	return status.Error(codes.Internal, fmt.Sprintf("unexpected error = %s", err))
}
