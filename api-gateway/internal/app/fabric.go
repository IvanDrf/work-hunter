package app

import (
	"time"

	"github.com/IvanDrf/work-hunter/api-gateway/internal/config"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/domain/ports"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/infrastructure/clients"
	"github.com/IvanDrf/work-hunter/api-gateway/internal/interfaces/http"
)

type fabric struct {
	config *config.Config
}

func NewFabric(config *config.Config) *fabric {
	return &fabric{
		config: config,
	}
}

func (f *fabric) NewApp() *App {
	auth := clients.NewAuthClient(f.config.Auth.Host, f.config.Auth.Port, f.config.App.Retries)
	vacancy := clients.NewVacancyClient(f.config.Vacancy.Host, f.config.Vacancy.Port, f.config.App.Retries)

	handlers := f.newHandlers(auth, vacancy, f.config.App.RequestTime)
	middleware := f.newAuthMiddleware(auth, f.config.App.RequestTime)

	server := f.newHttpServer(f.config.App.Host, f.config.App.Port, handlers, middleware)
	return &App{
		server: server,
	}
}

func (f *fabric) newHandlers(authClient ports.AuthClient, vacancyClient ports.VacancyClient, requestTime time.Duration) *http.Handlers {
	return http.NewHandlers(authClient, vacancyClient, requestTime)
}

func (f *fabric) newAuthMiddleware(authClient ports.AuthClient, requestTime time.Duration) *http.AuthMiddleware {
	return http.NewAuthMiddleware(authClient, requestTime)
}

func (f *fabric) newHttpServer(host string, port int, handlers *http.Handlers, middleware *http.AuthMiddleware) *http.Server {
	return http.NewServer(host, port, handlers, middleware)
}
