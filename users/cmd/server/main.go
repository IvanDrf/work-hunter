package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/config"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/database"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/logger"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/migration/migrator"
)

func main() {
	// init config
	config := config.MustLoad()

	// init logger
	log := logger.New(&config.Logger)
	log.Info("Starting user service", slog.Int("port", config.Server.Port), slog.String("log_level", config.Logger.Level))

	// connect to db
	db, err := database.NewPostgresConnection(config.Database, log)
	if err != nil {
		log.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// migrations
	if err := migrator.MigrateUp(config.Database.DSN()); err != nil {
		log.Error("failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	// TODO: init repo
	// TODO: init service
	// TODO: init GRPC server

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// TODO: start gRPC server

	<-ctx.Done()
	log.Info("Shutting down gracefully...")

	stutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: gRPC server gs

	<-stutdownCtx.Done()
	log.Info("Shutdown completed")

}
