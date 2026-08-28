package api

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/content/internal/app"
	"github.com/IvanDrf/work-hunter/content/internal/config"
	"github.com/IvanDrf/work-hunter/content/internal/infrastructure/adapters/logger"
)

func main() {
	cfg := config.LoadFromEnv()
	log := logger.SetupLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", "error", err.Error())
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		log.Error("application stopped with error", "error", err.Error())
		os.Exit(1)
	}
}
