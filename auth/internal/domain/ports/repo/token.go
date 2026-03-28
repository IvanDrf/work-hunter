package repo

import (
	"context"
	"time"
)

type TokenRepo interface {
	CreateToken(ctx context.Context, email string, token string, ttl time.Duration) error
	FindEmailByToken(ctx context.Context, token string) string

	DeleteToken(ctx context.Context, token string) error

	Close()
}
