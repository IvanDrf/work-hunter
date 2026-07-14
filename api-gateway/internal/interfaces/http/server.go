package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type server struct {
	handlers *handlers
	server   *http.Server
}

func NewServer(host string, port int, handlers *handlers) *server {
	return &server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
		},
		handlers: handlers,
	}
}

func (s *server) registerRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", s.handlers.RegisterUser)

	s.server.Handler = mux
}

func (s *server) Start() {
	l := slog.With(slog.String("server", "http"))

	s.registerRoutes()
	s.checkServicesHealth()

	l.Info("Starting server", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatalf("can't start http server, error=%s", err)
	}
}

func (s *server) checkServicesHealth() {
	const healthCheckTime = time.Second

	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTime)
	defer cancel()

	s.handlers.checkClientsHealth(ctx)
}

func (s *server) Close(ctx context.Context) {
	slog.InfoContext(ctx, "Stopping http server on", slog.String("address", s.server.Addr))

	s.server.Shutdown(ctx)
	s.handlers.close()
}
