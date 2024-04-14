package main

import (
	"avito2/internal/env/config"
	"avito2/internal/handler"
	"avito2/internal/infrastructure"
	"avito2/internal/repository"
	"avito2/internal/service"

	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {

	app := fiber.New()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	pgConfig, err := config.GetDBConfig(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("got bd config")

	db, err := infrastructure.SetUpPostgresDatabase(ctx, pgConfig)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Info("connected to db")

	authConfig, err := config.GetAuthConfig(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	redisconfig, err := config.GetRedisConfig(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	cache, err := infrastructure.SetUpRedisCache(ctx, redisconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	bannerRepository := repository.NewBannerRepository(db)
	bannerService := service.NewBannerService(bannerRepository, cache)
	bannerHandler := handler.NewBannerHandler(bannerService)

	bannerHandler.InitRoutes(app, authConfig)
	app.Listen(":3000")
}
