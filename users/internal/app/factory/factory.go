package factory

import "github.com/IvanDrf/workk-hunter/pkg/users/internal/config"

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}
