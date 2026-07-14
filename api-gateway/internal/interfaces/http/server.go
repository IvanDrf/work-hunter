package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
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
	mux.HandleFunc("POST /api/login", s.handlers.LoginUser)

	s.server.Handler = mux
}

func (s *server) Start(ctx context.Context) {
	l := slog.With(slog.String("server", "http"))

	s.registerRoutes()
	go s.checkServicesHealth(ctx)

	l.Info("Starting server", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start http server, error=%s", err)
	}
}

func (s *server) checkServicesHealth(ctx context.Context) {
	s.handlers.checkClientsHealth(ctx)
}

func (s *server) Close(ctx context.Context) {
	slog.InfoContext(ctx, "Stopping http server on", slog.String("address", s.server.Addr))

	s.server.Shutdown(ctx)
	s.handlers.close()
}
