package security

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

const JWT_CONTEXT_KEY = "userInfo"

func GetJwtTokenHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(viper.GetString("JWT_SECRET"))},
		ContextKey: JWT_CONTEXT_KEY,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
	})
}

func CheckTokenHeaderExist(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	} else {
		return ctx.Next()
	}
}
