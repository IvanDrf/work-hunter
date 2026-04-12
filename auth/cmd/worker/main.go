package main

// Email worker is used to send verification messages on email

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/auth/internal/app/factory"
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/adapters"
)

func main() {
	cfg := config.LoadFromYAML()
	adapters.InitLogger(&cfg.App)

	worker := factory.NewFactory(cfg).NewEmailWorker()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Start(ctx)
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()

	worker.Stop()
}
