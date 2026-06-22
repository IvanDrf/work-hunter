package redis

import (
	"context"
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_DSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("can't connect to redis server %s", err)
	}

	return client
}
