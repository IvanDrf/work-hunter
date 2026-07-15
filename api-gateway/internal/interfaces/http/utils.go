package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
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
	w.Header().Add("Content-type", "applications/json")

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

		case models.ErrCodeInvalidCookie:
			w.WriteHeader(http.StatusBadRequest)

		case models.ErrCodeAccess:
			w.WriteHeader(http.StatusForbidden)

		case models.ErrCodeNotFound:
			w.WriteHeader(http.StatusNotFound)
		}

		json.NewEncoder(w).Encode(e)

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

func createCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func getCookie(ctx context.Context, w http.ResponseWriter, r *http.Request, name string) (*http.Cookie, error) {
	log := adapters.GetLogger(ctx)

	cookie, err := r.Cookie(name)
	e := models.Error{
		Code: models.ErrCodeInvalidCookie,
	}

	if err != nil {
		e.Message = err.Error()
	} else if cookie.Valid() != nil {
		e.Message = cookie.Valid().Error()
	} else {
		return cookie, nil
	}

	log.ErrorContext(ctx, "invalid cookie", slog.String("error", e.Message))

	w.Header().Add("Content-type", "applications/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(e)
	return nil, e
}

type validator interface {
	IsValid() bool
}

func validateModel(ctx context.Context, w http.ResponseWriter, model validator, errorMessage string) error {
	if model.IsValid() {
		return nil
	}

	log := adapters.GetLogger(ctx)

	log.InfoContext(ctx, errorMessage, slog.String("error", errorMessage))

	err := models.Error{
		Message: errorMessage,
		Code:    models.ErrCodeUnprocessableEntity,
	}

	w.Header().Add("Content-type", "applications/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(err)

	return err
}
