package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/auth/internal/app/factory"
	"github.com/IvanDrf/work-hunter/auth/internal/config"
)

func main() {
	cfg := config.LoadFromYAML()

	worker := factory.NewFactory(cfg).NewEmailWorker()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Start(ctx)
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM)
	<-stop
	cancel()

	worker.Stop()
}
