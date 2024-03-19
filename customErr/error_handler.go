package customErr

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ErrorNotFound       = fiber.NewError(fiber.StatusNotFound, "not found")
	ErrorBadRequest     = fiber.NewError(fiber.StatusBadRequest, "bad request")
	ErrorUnauthorized   = fiber.NewError(fiber.StatusUnauthorized, "unauthorizex")
	ErrorConflict       = fiber.NewError(fiber.StatusConflict, "value already exists")
	ErrorInternalServer = fiber.NewError(fiber.StatusInternalServerError, "internal error")
)

func NotFoundHandler(ctx *fiber.Ctx) error {
	return ctx.Status(404).JSON(fiber.Map{
		"message": "Not Found",
	})
}
