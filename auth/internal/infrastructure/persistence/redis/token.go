package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepo struct {
	client *redis.Client
}

func NewTokenRepo(client *redis.Client) *TokenRepo {
	return &TokenRepo{
		client: client,
	}
}

func (t *TokenRepo) CreateToken(ctx context.Context, email string, token string, ttl time.Duration) error {
	return t.client.Set(ctx, token, email, ttl).Err()
}

func (t *TokenRepo) FindEmailByToken(ctx context.Context, token string) string {
	email := t.client.Get(ctx, token).Val()

	return email
}

func (t *TokenRepo) DeleteToken(ctx context.Context, token string) error {
	return t.client.Del(ctx, token).Err()
}

func (t *TokenRepo) Close() {
	t.client.Close()
}
