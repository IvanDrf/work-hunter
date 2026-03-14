package main

import (
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/config"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/logger"
)

func main() {
	config := config.MustLoad()

	log := logger.New(&config.Logger)

	log.Info("starting server ...")
	// TODO: init db ...
	// TODO: init server ...
}
