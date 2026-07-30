package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/IvanDrf/work-hunter/users/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	db *sqlx.DB
}

// create connection to postgres database
func NewPostgresConnection(cfg config.DBConfig) (*PostgresConnection, error) {
	db, err := sqlx.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresConnection{
		db: db,
	}, nil
}

func NewPostgresConnectionTest(db *sqlx.DB) *PostgresConnection {
	return &PostgresConnection{
		db: db,
	}
}

func (c *PostgresConnection) Close() {
	c.db.Close()
}
