package config

type AppConfig struct {
	LoggerLevel string `env:"LOGGER_LEVEL"`

	AppHost string `env:"APP_HOST"`
	AppPort int    `env:"APP_PORT"`
}
