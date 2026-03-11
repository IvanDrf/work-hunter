package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("can't connect to postgres database: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("can't ping postgres database: %s", err)
	}

	return db
}
