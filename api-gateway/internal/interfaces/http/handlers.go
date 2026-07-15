package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/ports"
)

type Handlers struct {
	authClient ports.AuthClient

	requestTime time.Duration
}

func NewHandlers(authClient ports.AuthClient, requestTime time.Duration) *Handlers {
	return &Handlers{
		authClient:  authClient,
		requestTime: requestTime,
	}
}

func (h *Handlers) close() {
	h.authClient.Close()
}

func (h *Handlers) checkClientsHealth(ctx context.Context) {
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

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "AVAILABLE",
	})
}
