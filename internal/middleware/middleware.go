package middleware

import (
	"avito2/internal/env/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func AdminAuth(config *config.AuthConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.AdminSecretKey),
		TokenLookup:  "header:token",
		ErrorHandler: customErrorHandler,
		ContextKey:   "user",
	})
}

func GeneralAuth(config *config.AuthConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKeys: map[string]interface{}{
			"user_secret":  []byte(config.UserSecretKey),
			"admin_secret": []byte(config.AdminSecretKey),
		},
		TokenLookup:  "header:token",
		ErrorHandler: customErrorHandler,
		ContextKey:   "user",
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Пользователь не авторизован",
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Внутренняя ошибка сервера",
	})
}
