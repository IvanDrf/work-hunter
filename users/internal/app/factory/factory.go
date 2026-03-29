package factory

import "github.com/IvanDrf/work-hunter/users/internal/config"

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}
