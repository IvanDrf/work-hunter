package config

type EmailConfig struct {
	EmailHost string `env:"EMAIL_HOST_ADDR"`

	EmailUsername string `env:"EMAIL_USERNAME"`
	EmailPassword string `env:"EMAIL_PASSWORD"`
}
