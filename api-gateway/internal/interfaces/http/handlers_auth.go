package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "register"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	if status, err := validateHeaders(r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Add("Content-type", "applications/json")
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid body request",
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	if err := validateModel(ctx, w, user, "email or password is empty"); err != nil {
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

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "login"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

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

	if err := validateModel(ctx, w, user, "email or password is empty"); err != nil {
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

func (h *Handlers) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "change-password"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	if status, err := validateHeaders(r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	password := &models.Password{}
	if err := json.NewDecoder(r.Body).Decode(password); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid body request",
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	if err := validateModel(ctx, w, password, "old or new password is empty"); err != nil {
		log.InfoContext(ctx, "invalid user content in request body", slog.String("error", err.Error()))
		return
	}

	access, err := getCookie(ctx, w, r, "access")
	if err != nil {
		log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
		return
	}

	if err := h.authClient.SendChangePasswordRequest(ctx, access.Value, password.Old, password.New); err != nil {
		handleResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "refresh-tokens"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	refresh, err := getCookie(ctx, w, r, "refresh")
	if err != nil {
		log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
		return
	}

	tokens, err := h.authClient.SendRefreshTokensRequest(ctx, refresh.Value)
	if err != nil {
		handleResponseError(w, err)
		return
	}

	access, refresh := setCookie("access", tokens.Access), setCookie("refresh", tokens.Refresh)
	http.SetCookie(w, access)
	http.SetCookie(w, refresh)

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "delete-user"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	if status, err := validateHeaders(r); err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	password := &struct {
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(password); err != nil {
		log.InfoContext(ctx, "can't parse requests's body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Error{
			Message: "invalid body request",
			Code:    models.ErrCodeUnprocessableEntity,
		})
		return
	}
	defer r.Body.Close()

	access, err := getCookie(ctx, w, r, "access")
	if err != nil {
		log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
		return
	}

	if err := h.authClient.SendDeleteUserRequest(ctx, access.Value, password.Password); err != nil {
		handleResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTime)
	defer cancel()

	log := slog.With(slog.String("handlers", "send-verfication-email"))
	log.InfoContext(ctx, "request")

	ctx = adapters.InsertLogger(ctx, log)

	access, err := getCookie(ctx, w, r, "access")
	if err != nil {
		log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
		return
	}

	if err := h.authClient.SendVerificationEmailRequest(ctx, access.Value); err != nil {
		handleResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
