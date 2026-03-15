package app

import (
	"fmt"
	"log"
	"net"

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
		handlers: newFactory(cfg).NewHandlers(),
		server:   grpc.NewServer(),
	}

	reflection.Register(app.server)
	auth_api.RegisterAuthServer(app.server, app.handlers)
	return app
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.Port))
	if err != nil {
		log.Fatalf("can't start AUTH service: %s", err)
	}

	if err := a.server.Serve(l); err != nil {
		log.Fatalf("can't start AUTH servie: %s", err)
	}
}

func (a *App) Stop() {
	a.server.GracefulStop()
	a.handlers.Close()
}
