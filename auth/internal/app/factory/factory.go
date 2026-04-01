package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
)

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}

func (f *Factory) NewHandlers() *handlers.Handler {
	// repos
	userRepo, tokenRepo := f.newRepos()

	// services
	auth, verif := f.newServices(userRepo, tokenRepo)

	return handlers.NewHandler(auth, verif)
}
