package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IvanDrf/work-hunter/users/internal/app"
	"github.com/IvanDrf/work-hunter/users/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app := app.NewApp(cfg)

	go app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop()
}
