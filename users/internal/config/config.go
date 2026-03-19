package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

// Service configuration
type Config struct {
	Server   ServerConfig `yaml:"server"`
	Logger   LoggerConfig `yaml:"logger"`
	Database DBConfig     `yaml:"database"`
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

// Database config PostgreSQL
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`

	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
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
