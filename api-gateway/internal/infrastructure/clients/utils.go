package clients

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleResponseError(err error) error {
	st, _ := status.FromError(err)
	st.Message()

	switch status.Code(err) {
	case codes.AlreadyExists:
		return models.Error{
			Message: fmt.Sprintf("already exists, error=%s", st.Message()),
			Code:    models.ErrCodeAlreadyExists,
		}

	case codes.Internal, codes.Unavailable:
		return models.Error{
			Message: fmt.Sprintf("internal error, error=%s", st.Message()),
			Code:    models.ErrCodeInternal,
		}

	case codes.InvalidArgument:
		return models.Error{
			Message: fmt.Sprintf("invalid argument, error=%s", st.Message()),
			Code:    models.ErrCodeInvalidArgument,
		}

	case codes.NotFound:
		return models.Error{
			Message: fmt.Sprintf("not found, error=%s", st.Message()),
			Code:    models.ErrCodeNotFound,
		}

	default:
		return models.Error{
			Message: fmt.Sprintf("unexpected error=%s", st.Message()),
			Code:    models.ErrCodeInternal,
		}
	}
}

const baseDelay = time.Duration(0.2 * float64(time.Second))

func retry(ctx context.Context, retries int, fn func() (any, error)) (any, error) {
	var (
		err   error
		res   any
		delay = baseDelay
	)

	retryCodes := map[codes.Code]bool{
		codes.Internal:    true,
		codes.Unavailable: true,
	}

	log := adapters.GetLogger(ctx)

	for attempt := range retries + 1 {
		if res, err = fn(); err != nil && retryCodes[status.Code(err)] {
			log.ErrorContext(ctx, "failed attempt", slog.Int("attempt", attempt+1))

			time.Sleep(delay)
			delay = baseDelay * time.Duration(attempt+1)
		} else if err != nil {
			return nil, handleResponseError(err)
		} else {
			return res, nil
		}
	}

	return res, handleResponseError(err)
}
