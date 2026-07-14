package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
)

func validateHeaders(r *http.Request) (int, error) {
	if r.Header.Get("Content-type") != "application/json" {
		return 0, models.Error{
			Message: fmt.Sprintf("content type is not application/json, type=%s", r.Header.Get("Content-type")),
			Code:    models.ErrCodeUnsupportedMediaType,
		}
	}

	return http.StatusUnsupportedMediaType, nil
}

func handleResponseError(w http.ResponseWriter, err error) {
	var e models.Error
	if errors.As(err, &e) {
		switch e.Code {
		case models.ErrCodeAlreadyExists:
			w.WriteHeader(http.StatusConflict)

		case models.ErrCodeInvalidArgument:
			w.WriteHeader(http.StatusBadRequest)

		case models.ErrCodeUnprocessableEntity:
			w.WriteHeader(http.StatusUnprocessableEntity)

		case models.ErrCodeInternal:
			w.WriteHeader(http.StatusInternalServerError)

		case models.ErrCodeUnsupportedMediaType:
			w.WriteHeader(http.StatusUnsupportedMediaType)
		}

		json.NewEncoder(w).Encode(e)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

func setCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func validateUser(ctx context.Context, w http.ResponseWriter, user *models.User, log *slog.Logger) error {
	if user.IsUserValid() {
		return nil
	}

	log.InfoContext(ctx, "invalid user content in request", slog.String("error", "email or password is empty"))
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(models.Error{
		Message: "invalid body request, email or password is empty",
		Code:    models.ErrCodeUnprocessableEntity,
	})

	return models.Error{
		Message: "invalid user content in request body",
		Code:    models.ErrCodeInvalidArgument,
	}
}
