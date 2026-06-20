package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/auth/internal/app"
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/adapters"
)

func main() {
	cfg := config.LoadFromENV()
	adapters.InitLogger(cfg)

	app := app.NewApp(cfg)

	go app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop()
}
