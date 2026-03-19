package factory

import (
	"database/sql"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/postgres"
)

func (f *Factory) newRepos() (repo.UserRepo, repo.TokenRepo) {
	db := postgres.Connect(&f.cfg.Database)

	return f.newUserRepo(db), f.newTokenRepo(db)
}

func (f *Factory) newUserRepo(db *sql.DB) repo.UserRepo {
	return postgres.NewUserRepo(db)
}

func (f *Factory) newTokenRepo(db *sql.DB) repo.TokenRepo {
	return postgres.NewTokenRepo(db)
}
