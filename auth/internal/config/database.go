package config

import "fmt"

type PostgreConfig struct {
	PostgresHost string `env:"POSTGRES_HOST"`
	PostgresPort uint16 `env:"POSTGRES_PORT"`

	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`

	PostgresDatabase string `env:"POSTGRES_DB_NAME"`
}

func (d *PostgreConfig) POSTGRES_DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.PostgresUsername, d.PostgresPassword, d.PostgresHost, d.PostgresPort, d.PostgresDatabase,
	)
}

type RedisConfig struct {
	RedisHost string `env:"REDIS_HOST"`
	RedisPort int    `env:"REDIS_PORT"`

	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDatabase int    `env:"REDIS_DATABASE"`
}

func (r *RedisConfig) REDIS_DSN() string {
	return fmt.Sprintf("%s:%d", r.RedisHost, r.RedisPort)
}
