package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/auth/pkg"

	s "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
)

func (f *Factory) newServices(userRepo repo.UserRepo, tokenRepo repo.TokenRepo) (service.AuthService, service.VerificationService) {
	jwter := f.newJwter()
	producer := f.newProducer()

	return f.newAuthService(userRepo, jwter), f.newVerificationService(producer, userRepo, tokenRepo, jwter)
}

func (f *Factory) newAuthService(userRepo repo.UserRepo, jwter jwt.Jwter) service.AuthService {
	return s.NewAuthService(
		userRepo,
		jwter,
	)
}

func (f *Factory) newVerificationService(producer service.EmailProducer, userRepo repo.UserRepo, tokenRepo repo.TokenRepo, jwter jwt.Jwter) service.VerificationService {
	return s.NewVerificationService(producer, userRepo, tokenRepo, jwter)
}

func (f *Factory) newJwter() jwt.Jwter {
	return pkg.NewJwt(f.cfg.Jwt.Secret, f.cfg.Jwt.AccessTime, f.cfg.Jwt.RefreshTime)
}
