package factory

import (
	repoPort "github.com/IvanDrf/work-hunter/users/internal/domain/ports/repo"
	servicePort "github.com/IvanDrf/work-hunter/users/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/users/internal/infrastructure/service"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
)

func (f *Factory) newServices(repo repoPort.UserRepository, log *logger.Logger) servicePort.UserService {
	return service.NewUserService(repo, log)
}
