package customErr

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrorNotFound     = errors.New("not found")
	ErrorBadRequest   = errors.New("bad request")
	ErrorUnauthorized = errors.New("Unauthorized")
	ErrorConflict     = errors.New("value is already exists")
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	switch {
	case errors.Is(err, ErrorNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, ErrorBadRequest):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, ErrorUnauthorized):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrorConflict):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return fiber.NewError(code, "Oops!, Something Error - Message : "+err.Error())
	}
}
