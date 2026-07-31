package config

import (
	"fmt"
	"net"
)

type PostgreConfig struct {
	PostgresHost string `env:"POSTGRES_HOST"`
	PostgresPort uint16 `env:"POSTGRES_PORT"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`

	PostgresDatabase string `env:"POSTGRES_DB"`
}

func (d *PostgreConfig) POSTGRES_DSN() string {
	return fmt.Sprintf( //nolint:nosprintfhostport
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.PostgresUser, d.PostgresPassword, d.PostgresHost, d.PostgresPort, d.PostgresDatabase,
	)
}

type RedisConfig struct {
	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`

	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDatabase int    `env:"REDIS_DATABASE"`
}

func (r *RedisConfig) REDIS_DSN() string {
	return net.JoinHostPort(r.RedisHost, r.RedisPort)
}
