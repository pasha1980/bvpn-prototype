package http

import (
	"github.com/gofiber/fiber/v2"
)

func errorHandler(ctx *fiber.Ctx, err error) error {

	var statusCode int
	var message string
	var errorCode int

	switch err.(type) {
	case *fiber.Error:
		statusCode = err.(*fiber.Error).Code
		errorCode = 0
		message = err.(*fiber.Error).Message
		break
	default:
		statusCode = 500
		errorCode = 1
		message = err.Error()
		break
	}

	ctx.Status(statusCode)
	return ctx.JSON(map[string]any{
		"code":  errorCode,
		"error": message,
	})
}
