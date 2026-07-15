package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

type Server struct {
	handlers   *Handlers
	middleware *AuthMiddleware

	server *http.Server
}

func NewServer(host string, port int, handlers *Handlers, middleware *AuthMiddleware) *Server {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
		},
		handlers:   handlers,
		middleware: middleware,
	}
}

func (s *Server) Start(ctx context.Context) {
	l := slog.With(slog.String("server", "http"))

	s.registerRoutes()
	go s.checkServicesHealth(ctx)

	l.Info("Starting server", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start http server, error=%s", err)
	}
}

func (s *Server) checkServicesHealth(ctx context.Context) {
	s.handlers.checkClientsHealth(ctx)
}

func (s *Server) Close(ctx context.Context) {
	slog.InfoContext(ctx, "Stopping http server on", slog.String("address", s.server.Addr))

	s.server.Shutdown(ctx)
	s.handlers.close()
}
