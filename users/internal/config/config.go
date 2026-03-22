package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v3"
)

// Service configuration
type Config struct {
	Server   ServerConfig `yaml:"server"`
	Logger   LoggerConfig `yaml:"logger"`
	Database DBConfig     `yaml:"database"`
}

func MustLoad() *Config {
	filepath := os.Getenv("CONFIG_PATH")

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("failed to read config file %s: %v", filepath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to unmarshal data: %v", err)
	}

	return &cfg
}
