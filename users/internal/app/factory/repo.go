package factory

import (
	repository "github.com/IvanDrf/workk-hunter/pkg/users/internal/domain/ports/repo"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/infrastructure/persistence/postgres"
)

func (f *Factory) newRepos() (repository.UserRepository, error) {
	conn, err := postgres.NewPostgresConnection(f.cfg.Database)
	if err != nil {
		return nil, err
	}

	return postgres.NewUserRepository(conn), nil
}
