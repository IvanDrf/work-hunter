package config

type DatabaseConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Name string `yaml:"db_name"`
}
