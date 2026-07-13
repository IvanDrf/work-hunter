package main

import (
	"github.com/IvanDrf/work-hunter/api-gateway/internal/config"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/clients"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

func main() {
	config := config.LoadFromENV()

	s := http.NewServer(config.App.Host, config.App.Port, http.NewHandlers(
		clients.NewAuthClient(config.Auth.Host, config.Auth.Port, config.App.Retries), config.App.RequestTime,
	))

	s.Start()

}
