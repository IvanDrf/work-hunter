package config

import "time"

type JwtConfig struct {
	JWTSecret string `env:"JWT_SECRET"`

	JWTAccessTime  time.Duration `env:"JWT_ACCESS_TIME"`
	JWTRefreshTime time.Duration `env:"JWT_REFRESH_TIME"`
}
