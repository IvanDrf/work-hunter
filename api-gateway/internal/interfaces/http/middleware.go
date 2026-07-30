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

type payload string

const payloadKey payload = "payload"

func (m *AuthMiddleware) RegistredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.requestTime)
		defer cancel()

		log := slog.With(slog.String("middleware", "RegistredMiddleware"))
		ctx = adapters.InsertLogger(ctx, log)

		payload, err := m.getPyload(ctx, w, r)
		if err != nil {
			log.InfoContext(ctx, "invalid token payload", slog.String("error", err.Error()))
			handleResponseError(w, err)
			return
		}

		ctx = context.WithValue(ctx, payloadKey, &models.UserInfo{
			Role:        payload.Role,
			UserID:      payload.ID.String(),
			Verificated: payload.Verificated,
		})
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) ProbablyUnregistredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.requestTime)
		defer cancel()

		log := slog.With(slog.String("middleware", "ProbablyUnregistredMiddleware"))
		ctx = adapters.InsertLogger(ctx, log)

		payload, err := m.getPyload(ctx, w, r)
		if err == nil {
			next(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(ctx, payloadKey, &models.UserInfo{
			Role:        payload.Role,
			UserID:      payload.ID.String(),
			Verificated: payload.Verificated,
		})

		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.requestTime)
		defer cancel()

		log := slog.With(slog.String("middleware", "AdminMiddleware"))
		ctx = adapters.InsertLogger(ctx, log)

		payload, err := m.getPyload(ctx, w, r)
		if err != nil {
			log.InfoContext(ctx, "invalid token payload", slog.String("error", err.Error()))
			handleResponseError(w, err)
			return
		}

		if payload.Role != models.ADMIN {
			handleResponseError(w, models.Error{
				Message: "only admin user allowed",
				Code:    models.ErrCodeAccess,
			})
			return
		}

		ctx = context.WithValue(ctx, payloadKey, &models.UserInfo{
			Role:        payload.Role,
			UserID:      payload.ID.String(),
			Verificated: payload.Verificated,
		})
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) getPyload(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.TokenPayload, error) {
	log := adapters.GetLogger(ctx)

	access, err := getCookie(ctx, w, r, "access")
	if err != nil {
		log.InfoContext(ctx, "invalid cookie in request", slog.String("error", err.Error()))
		return nil, err
	}

	payload, err := m.authClient.SendIsTokenValidRequest(ctx, access.Value)
	return payload, err
}
