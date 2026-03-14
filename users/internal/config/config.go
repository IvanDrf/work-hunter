package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v3"
)

// Service configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

// Server configuration
type ServerConfig struct {
	Port int `yaml:"port"`
}

// Logger configuration
type LoggerConfig struct {
	Level      string   `yaml:"level"`
	Format     string   `yaml:"format"`
	OutputPath []string `yaml:"output_path"`
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
