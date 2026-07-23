package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/app/factory"
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg *config.Config

	handlers *handlers.Handler
	server   *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	app := &App{
		cfg:      cfg,
		handlers: factory.NewFactory(cfg).NewHandlers(),
		server:   grpc.NewServer(),
	}

	reflection.Register(app.server)
	auth_api.RegisterAuthServer(app.server, app.handlers)
	return app
}

func (a *App) Run() {
	const connectionTime = 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), connectionTime)
	defer cancel()

	config := &net.ListenConfig{}

	l, err := config.Listen(ctx, "tcp", fmt.Sprintf("%s:%d", a.cfg.AppHost, a.cfg.AppPort))
	if err != nil {
		cancel()
		log.Printf("can't start AUTH service: %s", err)
		return
	}

	slog.Info("Starting AUTH service", slog.String("host", a.cfg.AppHost), slog.Int("port", a.cfg.AppPort))

	if err := a.server.Serve(l); err != nil {
		cancel()
		log.Printf("can't start AUTH servie: %s", err)
		return
	}
}

func (a *App) Stop() {
	slog.Info("Stopping AUTH service", slog.String("host", a.cfg.AppHost), slog.Int("port", a.cfg.AppPort))

	a.server.GracefulStop()
	a.handlers.Close()
}
