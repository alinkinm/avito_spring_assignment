package infrastructure

import (
	cache "avito2/internal/cache"
	"avito2/internal/env/config"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetUpRedisCache(ctx context.Context, config *config.RedisConfig) (*cache.MyCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	return &cache.MyCache{Client: client}, nil
}
