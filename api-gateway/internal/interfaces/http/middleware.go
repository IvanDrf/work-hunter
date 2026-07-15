package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/models"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/ports"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/adapters"
)

type AuthMiddleware struct {
	authClient  ports.AuthClient
	requestTime time.Duration
}

func NewAuthMiddleware(authClient ports.AuthClient, requestTime time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:  authClient,
		requestTime: requestTime,
	}
}

func (m *AuthMiddleware) RegistredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.requestTime)
		defer cancel()

		log := slog.With(slog.String("middleware", "RegistredMiddleware"))
		ctx = adapters.InsertLogger(ctx, log)

		access, err := getCookie(ctx, w, r, "access")
		if err != nil {
			log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
			return
		}

		payload, err := m.authClient.SendIsTokenValidRequest(ctx, access.Value)
		if err != nil {
			handleResponseError(w, err)
			return
		}

		ctx = context.WithValue(ctx, "payload", payload)
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.requestTime)
		defer cancel()

		log := slog.With(slog.String("middleware", "AdminMiddleware"))
		ctx = adapters.InsertLogger(ctx, log)

		access, err := getCookie(ctx, w, r, "access")
		if err != nil {
			log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
			return
		}

		payload, err := m.authClient.SendIsTokenValidRequest(ctx, access.Value)
		if err != nil {
			handleResponseError(w, err)
			return
		}

		if payload.Role != models.ADMIN {
			handleResponseError(w, models.Error{
				Message: "only admin user allowed",
				Code:    models.ErrCodeAccess,
			})
		}

		ctx = context.WithValue(ctx, "payload", payload)
		next(w, r.WithContext(ctx))
	}
}
