package http_in

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func errorHandler(ctx *fiber.Ctx, err error) error {

	var code int
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	} else {
		code = 400
	}

	ctx.Status(code)
	return ctx.JSON(map[string]string{
		"type": err.Error(),
	})
}
