package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key, value string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() Cache {
	c := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDDISADDR"),
		Username: "default",
		Password: os.Getenv("REDDISPWD"),
		DB:       0,
	})

	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %s", err)
	}

	return &redisCache{
		client: c,
	}
}

func (c *redisCache) Set(ctx context.Context, key, value string, exp time.Duration) error {
	err := c.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
