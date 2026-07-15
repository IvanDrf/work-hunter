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

	// Создаем sqlx.DB из sql.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Создаем PostgresConnection с мок-базой
	conn := &postgres.PostgresConnection{
		DB: sqlxDB,
	}

	repo := postgres.NewUserRepository(conn)
	return repo, mock
}
