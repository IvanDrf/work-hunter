package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/ports"
)

const (
	invalidBodyRequestMessage = "invalid body request"
)

type Handlers struct {
	authClient    ports.AuthClient
	vacancyClient ports.VacancyClient

	requestTime time.Duration
}

func NewHandlers(authClient ports.AuthClient, vacancyClient ports.VacancyClient, requestTime time.Duration) *Handlers {
	return &Handlers{
		authClient:    authClient,
		vacancyClient: vacancyClient,
		requestTime:   requestTime,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "AVAILABLE",
	})
}

func (h *Handlers) close() {
	h.authClient.Close()
}

func (h *Handlers) checkClientsHealth(ctx context.Context, healthCheckTime time.Duration) {
	defer slog.InfoContext(ctx, "stoppping clients health check")
	ticker := time.NewTicker(healthCheckTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.authClient.Health(ctx)
			h.vacancyClient.Health(ctx)
		case <-ctx.Done():
			return
		}
	}
}
