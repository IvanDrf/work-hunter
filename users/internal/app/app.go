package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	"github.com/IvanDrf/work-hunter/users/internal/app/factory"
	"github.com/IvanDrf/work-hunter/users/internal/config"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	cfg *config.Config
	log *logger.Logger

	handlers *handlers.Handler
	server   *grpc.Server
}

func NewApp(cfg *config.Config) *App {
	log := logger.New(&cfg.Logger)
	handlers, err := factory.NewFactory(cfg).NewHandler(log)
	if err != nil {
		log.Error("can`t build handlers", "error", err)
		os.Exit(1)
	}

	app := &App{
		cfg:      cfg,
		log:      log,
		handlers: handlers,
		server:   grpc.NewServer(),
	}

	reflection.Register(app.server)
	user_api.RegisterUserServer(app.server, app.handlers)
	return app
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.App.Host, a.cfg.App.Port))
	if err != nil {
		a.log.Error("can`t start USER service", "error", err)
		os.Exit(1)
	}

	a.log.Info("Starting USER service", slog.String("host", a.cfg.App.Host), slog.Int("port", a.cfg.App.Port))

	if err := a.server.Serve(l); err != nil {
		a.log.Error("can`t start USER service", "error", err)
		os.Exit(1)
	}
}

func (a *App) Stop() {
	a.log.Info("Stopping USER service", slog.String("host", a.cfg.App.Host), slog.Int("port", a.cfg.App.Port))

	a.server.GracefulStop()
	a.handlers.Close()
}
