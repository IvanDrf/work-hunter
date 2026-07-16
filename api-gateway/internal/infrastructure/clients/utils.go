package clients

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleResponseError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return models.Error{
			Message: fmt.Sprintf("unexpected error: %v", err),
			Code:    models.ErrCodeInternal,
		}
	}

	code := status.Code(err)
	msg := st.Message()

	groups := map[codes.Code]models.ErrCode{
		codes.AlreadyExists:      models.ErrCodeAlreadyExists,
		codes.InvalidArgument:    models.ErrCodeInvalidArgument,
		codes.OutOfRange:         models.ErrCodeInvalidArgument,
		codes.NotFound:           models.ErrCodeNotFound,
		codes.Unauthenticated:    models.ErrCodeAccess,
		codes.PermissionDenied:   models.ErrCodeAccess,
		codes.Canceled:           models.ErrCodeCanceled,
		codes.DeadlineExceeded:   models.ErrCodeDeadlineExceeded,
		codes.ResourceExhausted:  models.ErrCodeResourceExhausted,
		codes.FailedPrecondition: models.ErrCodeFailedPrecondition,
		codes.Aborted:            models.ErrCodeAborted,
		codes.Internal:           models.ErrCodeInternal,
		codes.Unavailable:        models.ErrCodeInternal,
		codes.Unknown:            models.ErrCodeInternal,
		codes.Unimplemented:      models.ErrCodeInternal,
		codes.DataLoss:           models.ErrCodeInternal,
	}

	errCode, exists := groups[code]
	label := "unexpected"
	if exists {
		label = strings.ReplaceAll(strings.ToLower(string(errCode)), "_", " ")
	}

	return models.Error{
		Message: fmt.Sprintf("%s: %s", label, msg),
		Code:    errCode,
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
