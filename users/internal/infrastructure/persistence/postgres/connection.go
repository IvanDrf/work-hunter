package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/IvanDrf/workk-hunter/pkg/users/internal/config"
	"github.com/IvanDrf/workk-hunter/pkg/users/internal/logger"
)

type PostgresConnection struct {
	db  *sqlx.DB
	log *logger.Logger
}

// create connection to postgres database
func NewPostgresConnection(cfg config.DBConfig, log *logger.Logger) (*PostgresConnection, error) {
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

	log.Info("Connected to database",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("dbname", cfg.DBName))

	return &PostgresConnection{
		db:  db,
		log: log,
	}, nil
}

func (c *PostgresConnection) GetDB() *sqlx.DB {
	return c.db
}

func (c *PostgresConnection) Close() error {
	return c.db.Close()
}

// execute function in transaction
func WithTransaction(ctx context.Context, db *sqlx.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rbErr, err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
