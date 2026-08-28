package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig
	S3Config
	RabbitMQConfig
}

func LoadFromEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't load config settings from .env, error=%s", err)
	}

	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("can't load config from env file, error=%s", err)
	}

	return cfg
}
