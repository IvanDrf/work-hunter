package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/auth/internal/app"
	"github.com/IvanDrf/work-hunter/auth/internal/config"
)

func main() {
	cfg := config.LoadFromYAML()

	app := app.NewApp(cfg)

	go app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM)

	<-stop
	app.Stop()
}
