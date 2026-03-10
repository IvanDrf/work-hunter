package app

import (
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/postgres"
	"github.com/IvanDrf/work-hunter/auth/internal/interfaces/grpc/handlers"
	"github.com/IvanDrf/work-hunter/auth/pkg"

	s "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
)

type Factory struct {
	cfg *config.Config
}

func newFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}

func (f *Factory) NewHandlers() *handlers.Handler {
	return handlers.NewHandler(f.newAuthService())
}

func (f *Factory) newAuthService() service.AuthService {
	return s.NewAuthService(
		f.newUserRepo(),
		f.newJwter(),
	)
}

func (f *Factory) newJwter() jwt.Jwter {
	return pkg.NewJwt(f.cfg.Jwt.Secret, f.cfg.Jwt.AccessTime, f.cfg.Jwt.RefreshTime)
}

func (f *Factory) newUserRepo() repo.UserRepo {
	return postgres.NewAuthRepo(&f.cfg.Database)
}
