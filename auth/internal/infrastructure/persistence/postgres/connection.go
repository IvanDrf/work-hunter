package postgres

import (
	"database/sql"
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.PostgreConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("can't connect to postgres database: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("can't ping postgres database: %s", err)
	}

	return db
}
