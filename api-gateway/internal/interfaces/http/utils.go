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
	if r.Header.Get("Content-Type") != "application/json" {
		return 0, models.Error{
			Message: fmt.Sprintf("content type is not application/json, type=%s", r.Header.Get("Content-Type")),
			Code:    models.ErrCodeUnsupportedMediaType,
		}
	}

	return http.StatusUnsupportedMediaType, nil
}

func handleResponseError(w http.ResponseWriter, err error) {
	var e models.Error
	w.Header().Add("Content-Type", "applications/json")

	statuses := map[models.ErrCode]int{
		models.ErrCodeAlreadyExists:        http.StatusConflict,
		models.ErrCodeInvalidArgument:      http.StatusBadRequest,
		models.ErrCodeUnprocessableEntity:  http.StatusUnprocessableEntity,
		models.ErrCodeInternal:             http.StatusInternalServerError,
		models.ErrCodeUnsupportedMediaType: http.StatusUnsupportedMediaType,
		models.ErrCodeInvalidCookie:        http.StatusBadRequest,
		models.ErrCodeAccess:               http.StatusForbidden,
		models.ErrCodeNotFound:             http.StatusNotFound,
	}

	if errors.As(err, &e) && statuses[e.Code] != 0 {
		w.WriteHeader(statuses[e.Code])
		json.NewEncoder(w).Encode(e)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

type cookieName string

const (
	Access  cookieName = "access"
	Refresh cookieName = "refresh"
)

func createCookie(name cookieName, value string) *http.Cookie {
	return &http.Cookie{
		Name:     string(name),
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func getCookie(ctx context.Context, w http.ResponseWriter, r *http.Request, name cookieName) (*http.Cookie, error) {
	log := adapters.GetLogger(ctx)

	cookie, err := r.Cookie(string(name))
	e := models.Error{
		Code: models.ErrCodeInvalidCookie,
	}

	switch {
	case err != nil:
		e.Message = err.Error()
	case cookie.Valid() != nil:
		e.Message = cookie.Valid().Error()
	default:
		return cookie, nil
	}

	log.ErrorContext(ctx, "invalid cookie", slog.String("error", e.Message))

	w.Header().Add("Content-Type", "applications/json")
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

	w.Header().Add("Content-Type", "applications/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(err)

	return err
}

func getUserInfo(ctx context.Context) (*models.UserInfo, error) {
	val := ctx.Value("user_info")

	userInfo, ok := val.(*models.UserInfo)
	if !ok {
		return nil, models.Error{
			Message: "can't get user info",
			Code:    models.ErrCodeInternal,
		}
	}

	return userInfo, nil
}
