package factory

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	s "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/service"
)

func (f *Factory) newSmtpEmailService() service.EmailService {
	return s.NewSmtpEmailService(f.cfg.Email.Host, f.cfg.Email.Username, f.cfg.Email.Password)
}
