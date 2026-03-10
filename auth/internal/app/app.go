package app

import (
	"fmt"
	"log"
	"net"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config

	server *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	app := &App{
		cfg:    cfg,
		server: grpc.NewServer(),
	}

	auth_api.RegisterAuthServer(app.server, newFactory(cfg).NewHandlers())
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
}
