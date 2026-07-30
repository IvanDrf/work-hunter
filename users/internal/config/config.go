package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Service configuration
type Config struct {
	App      AppConfig
	Logger   LoggerConfig
	Database DBConfig
}

func MustLoad() *Config {
	configPath := os.Getenv("ENV_FILE_PATH")
	if configPath == "" {
		configPath = ".env"
	}

	var cfg Config

	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("failed to read config file %s: %v", configPath, err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("failed to read env variables: %v", err)
		}
	}

	return &cfg
}
