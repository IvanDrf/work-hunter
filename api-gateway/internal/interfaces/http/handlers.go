package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/ports"
)

type handlers struct {
	authClient ports.AuthClient

	requestTime time.Duration
}

func NewHandlers(authClient ports.AuthClient, requestTime time.Duration) *handlers {
	return &handlers{
		authClient:  authClient,
		requestTime: requestTime,
	}
}

func (h *handlers) close() {
	h.authClient.Close()
}

func (h *handlers) checkClientsHealth(ctx context.Context) {
	defer slog.InfoContext(ctx, "stoppping clients health check")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.authClient.Health(ctx)
		case <-ctx.Done():
			return
		}
	}

}

func (h *handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "register"))
	log.InfoContext(ctx, "request")

	w.Header().Add("Content-type", "applications/json")
	if status, err := validateHeaders(r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid body request",
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	if err := validateUser(ctx, w, user, log); err != nil {
		log.InfoContext(ctx, "invalid user content in request body", slog.String("error", err.Error()))
		return
	}

	tokens, err := h.authClient.SendRegisterRequest(ctx, user.Email, user.Password, user.Role)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	access, refresh := setCookie("access", tokens.Access), setCookie("refresh", tokens.Refresh)
	http.SetCookie(w, access)
	http.SetCookie(w, refresh)

	w.WriteHeader(http.StatusNoContent)
}

func (h *handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "login"))
	log.InfoContext(ctx, "request")

	w.Header().Add("Content-type", "applications/json")
	if status, err := validateHeaders(r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid body request",
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	if err := validateUser(ctx, w, user, log); err != nil {
		log.InfoContext(ctx, "invalid user content in request body", slog.String("error", err.Error()))
		return
	}

	tokens, err := h.authClient.SendLoginRequest(ctx, user.Email, user.Password)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	access, refresh := setCookie("access", tokens.Access), setCookie("refresh", tokens.Refresh)
	http.SetCookie(w, access)
	http.SetCookie(w, refresh)

	w.WriteHeader(http.StatusNoContent)
}
