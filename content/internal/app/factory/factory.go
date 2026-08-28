package factory

import (
	"log/slog"

	"github.com/IvanDrf/work-hunter/content/internal/config"
)

type Factory struct {
	cfg *config.Config
	log *slog.Logger
}

func NewFactory(cfg *config.Config, log *slog.Logger) *Factory {
	return &Factory{
		cfg: cfg,
		log: log,
	}
}
