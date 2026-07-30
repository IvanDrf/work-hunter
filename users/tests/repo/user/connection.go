package user

import (
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanDrf/work-hunter/users/internal/infrastructure/persistence/postgres"
	"github.com/jmoiron/sqlx"
)

func connect() (*postgres.UserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	conn := postgres.NewPostgresConnectionTest(sqlxDB)

	repo := postgres.NewUserRepository(conn)
	return repo, mock
}
