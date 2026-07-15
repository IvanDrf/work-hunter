package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/app"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/config"
)

func main() {
	config := config.LoadFromENV()

	app := app.NewFabric(config).NewApp()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go app.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	app.Stop(ctx)
}
