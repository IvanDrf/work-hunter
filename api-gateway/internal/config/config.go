package config

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Host string `env:"APP_HOST"`
		Port int    `env:"APP_PORT"`

		LoggerLevel string `env:"LOGGER_LEVEL"`
		LoggerFile  string `env:"LOGGER_FILE"`

		RequestTime time.Duration `env:"REQUEST_TIME"`
		Retries     int           `env:"RETRIES"`

		ShutdownTime    time.Duration `env:"SHUTDOWN_TIME"`
		HealthCheckTime time.Duration `env:"HEALTH_CHECK_TIME"`
	}

	Auth struct {
		Host string `env:"AUTH_HOST"`
		Port int    `env:"AUTH_PORT"`
	}

	Applications struct {
		Host string `env:"APPLICATIONS_HOST"`
		Port int    `env:"APPLICATIONS_PORT"`
	}

	Vacancy struct {
		Host string `env:"VACANCY_HOST"`
		Port int    `env:"VACANCY_PORT"`
	}

	Validation struct {
		Host string `env:"VALIDATION_HOST"`
		Port int    `env:"VALIDATION_PORT"`
	}
}

func LoadFromENV() *Config {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't load .env config file, error=%s", err)
	}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("can't parse .env config file, error=%s", err)
	}

	return cfg
}
