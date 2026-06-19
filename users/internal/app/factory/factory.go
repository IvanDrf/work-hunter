package factory

import (
	"github.com/IvanDrf/work-hunter/users/internal/config"
	"github.com/IvanDrf/work-hunter/users/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
)

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}

func (f *Factory) NewHandler(log *logger.Logger) (*handlers.Handler, error) {
	userRepo, err := f.newRepos()
	if err != nil {
		return nil, err
	}

	service := f.newServices(userRepo, log)

	return handlers.NewHandler(service, log), nil
}
