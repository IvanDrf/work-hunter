package config

import "fmt"

type PostgreConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Name string `yaml:"db_name"`
}

func (d *PostgreConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.Username, d.Password, d.Host, d.Port, d.Name,
	)
}
