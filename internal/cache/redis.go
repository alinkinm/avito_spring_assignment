package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type MyCache struct {
	Client *redis.Client
}

func (c *MyCache) Get(key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *MyCache) Set(key, value string, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}
