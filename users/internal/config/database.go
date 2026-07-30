package config

import (
	"fmt"
	"time"
)

// Database config PostgreSQL
type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `env:"DB_NAME" env-default:"users_db"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`

	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" env-default:"1h"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}
