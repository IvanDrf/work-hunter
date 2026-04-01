package config

type EmailConfig struct {
	Host string `yaml:"host_addr"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
