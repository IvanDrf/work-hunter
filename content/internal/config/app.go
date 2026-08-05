package config

// app config
type AppConfig struct {
	Host   string `env:"APP_HOST"`
	Port   int    `env:"APP_PORT"`
	APIKey string `env:"API_KEY"`
	Env    string `env:"APP_ENV"`
}
