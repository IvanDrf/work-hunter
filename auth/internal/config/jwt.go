package config

import "time"

type JwtConfig struct {
	Secret string `yaml:"secret"`

	AccessTime  time.Duration `yaml:"access_time"`
	RefreshTime time.Duration `yaml:"refresh_time"`
}
