package config

type AppConfig struct {
	LoggerLevel string `yaml:"logger_level"`

	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
