package factory

import (
	"database/sql"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/postgres"
	r "github.com/IvanDrf/work-hunter/auth/internal/infrastructure/persistence/redis"
	"github.com/redis/go-redis/v9"
)

func (f *Factory) newRepos() (repo.UserRepo, repo.TokenRepo) {
	db := postgres.Connect(&f.cfg.Database)
	client := r.Connect(&f.cfg.Redis)

	return f.newUserRepo(db), f.newTokenRepo(client)
}

func (f *Factory) newUserRepo(db *sql.DB) repo.UserRepo {
	return postgres.NewUserRepo(db)
}

func (f *Factory) newTokenRepo(client *redis.Client) repo.TokenRepo {
	return r.NewTokenRepo(client)
}
