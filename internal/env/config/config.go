package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
)

type AuthConfig struct {
	UserSecretKey  string `env:"USER_SECRET_KEY"`
	AdminSecretKey string `env:"ADMIN_SECRET_KEY"`
}

type DBConfig struct {
	Name     string `env:"DB_NAME"`
	Password string `env:"DB_PASSWORD"`
	Db       string `env:"DB_DB"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
}

type RedisConfig struct {
	Address  string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASS"`
}

func GetDBConfig(ctx context.Context) (*DBConfig, error) {

	var c DBConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &c, nil
}

func GetRedisConfig(ctx context.Context) (*RedisConfig, error) {

	var c RedisConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &c, nil

}

func GetAuthConfig(ctx context.Context) (*AuthConfig, error) {

	var c AuthConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &c, nil

}
