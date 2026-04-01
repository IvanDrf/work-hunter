package factory

import (
	ports "github.com/IvanDrf/work-hunter/users/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/users/internal/infrastructure/persistence/postgres"
)

func (f *Factory) newRepos() (ports.UserRepository, error) {
	conn, err := postgres.NewPostgresConnection(f.cfg.Database)
	if err != nil {
		return nil, err
	}

	return postgres.NewUserRepository(conn), nil
}
