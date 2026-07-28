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
	userClient    ports.UserClient

	requestTime time.Duration
}

func NewHandlers(authClient ports.AuthClient, vacancyClient ports.VacancyClient, userClient ports.UserClient, requestTime time.Duration) *Handlers {
	return &Handlers{
		authClient:    authClient,
		vacancyClient: vacancyClient,
		userClient:    userClient,
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
	h.userClient.Close()
	h.vacancyClient.Close()
}

func (h *Handlers) checkClientsHealth(ctx context.Context, periodCheckTime time.Duration) {
	defer slog.InfoContext(ctx, "stoppping clients health check")
	ticker := time.NewTicker(periodCheckTime)
	defer ticker.Stop()

	time.Sleep(1 * time.Second)
	h.sendHealthChecks(ctx)

	for {
		select {
		case <-ticker.C:
			h.sendHealthChecks(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (h *Handlers) sendHealthChecks(ctx context.Context) {
	h.authClient.Health(ctx)
	h.userClient.Health(ctx)
	h.vacancyClient.Health(ctx)
}
