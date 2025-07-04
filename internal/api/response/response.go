package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func BadRequestResponse(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message": msg,
	})
}

func UnauthorizedResponse(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
		"message": msg,
	})
}

func NotFoundResponse(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusNotFound).JSON(&fiber.Map{
		"message": msg,
	})
}

func InternalServerError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
		"message": "server error",
		"error":   err.Error(),
	})
}

func SuccessResponse(ctx *fiber.Ctx, msg string, data any) error {
	if msg == "" {
		msg = "success response"
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}

func CreastedResponse(ctx *fiber.Ctx, msg string, data any) error {
	if msg == "" {
		msg = "data created"
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}
