package config

// Logger configuration
type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL" env-default:"info"`
	Format string `env:"LOG_FORMAT" env-default:"json"`
}
