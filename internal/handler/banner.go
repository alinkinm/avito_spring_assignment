package handler

import (
	core "avito2/internal/core"
	"avito2/internal/env/config"
	"context"

	"avito2/internal/middleware"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

const (
	contextKeyUser = "user"
)

type BannerService interface {
	GetUserBanner(ctx context.Context, req *core.BannerRequest, role string) (string, error)
	GetAllBanners(ctx context.Context, banner *core.BannerRequest2) ([]*core.Banner, error)
}

type BannerHandler struct {
	bannerService BannerService
}

func NewBannerHandler(service BannerService) *BannerHandler {
	return &BannerHandler{bannerService: service}
}

func (handler *BannerHandler) InitRoutes(app *fiber.App, config *config.AuthConfig) {
	app.Get("/user_banner", middleware.GeneralAuth(config), handler.GetUserBanner)
	app.Get("/banner", middleware.AdminAuth(config), handler.GetAllBanners)
}

func (handler *BannerHandler) GetUserBanner(ctx *fiber.Ctx) error {

	jwtPayload, ok := jwtPayloadFromRequest(ctx)
	if !ok {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	role := jwtPayload["role"].(string)

	var req core.BannerRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		ctx.Set("Content-Type", "application/json")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	banner, err := handler.bannerService.GetUserBanner(ctx.UserContext(), &req, role)
	if err != nil {
		if err == core.NewErrBannerDoesNotExist() {
			return ctx.Status(fiber.StatusNotFound).SendString("")
		}

		if err == core.NewErrInternalServerError() {
			ctx.Set("Content-Type", "application/json")
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	ctx.Set("Content-Type", "application/json")
	return ctx.Status(fiber.StatusOK).JSON(banner)

}

func (handler *BannerHandler) GetAllBanners(ctx *fiber.Ctx) error {

	var req core.BannerRequest2

	if err := ctx.QueryParser(&req); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		ctx.Set("Content-Type", "application/json")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil
}

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {

	jwtToken, ok := c.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		log.WithFields(log.Fields{
			"jwt_token_context_value": c.Context().Value(contextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.WithFields(log.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true

}
