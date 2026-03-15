package app

import (
	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	messaging "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/messaging/rabbitmq"
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
	// dependencies
	userRepo := f.newUserRepo()
	jwter := f.newJwter()
	emailProducer := f.newEmailProducer()

	// services
	auth := f.newAuthService(userRepo, jwter)
	verification := f.newVerificationService(emailProducer, userRepo, jwter)

	return handlers.NewHandler(auth, verification)
}

func (f *Factory) newAuthService(userRepo repo.UserRepo, jwter jwt.Jwter) service.AuthService {
	return s.NewAuthService(
		userRepo,
		jwter,
	)
}

func (f *Factory) newVerificationService(producer service.EmailProducer, userRepo repo.UserRepo, jwter jwt.Jwter) service.VerificationService {
	return s.NewVerificationService(producer, userRepo, jwter)
}

func (f *Factory) newJwter() jwt.Jwter {
	return pkg.NewJwt(f.cfg.Jwt.Secret, f.cfg.Jwt.AccessTime, f.cfg.Jwt.RefreshTime)
}

func (f *Factory) newUserRepo() repo.UserRepo {
	return postgres.NewAuthRepo(&f.cfg.Database)
}

func (f *Factory) newEmailProducer() service.EmailProducer {
	return messaging.NewRabbitMqProducer(&f.cfg.Broker)
}
