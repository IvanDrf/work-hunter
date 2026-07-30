package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.POSTGRES_DSN())
	if err != nil {
		log.Fatalf("can't connect to postgres database: %s", err)
	}

	const pingTime = 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), pingTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		cancel()
		log.Printf("can't ping postgres database: %s", err)
		return nil
	}

	return db
}
