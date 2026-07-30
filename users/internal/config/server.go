package config

// Server configuration
type AppConfig struct {
	Host string `env:"APP_HOST" env-default:"0.0.0.0"`
	Port int    `env:"APP_PORT" env-default:"8080"`
}
