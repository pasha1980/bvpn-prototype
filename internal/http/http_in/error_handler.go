package http_in

import (
	"bvpn-prototype/internal/protocols/protocol_error"
	"github.com/gofiber/fiber/v2"
)

func errorHandler(ctx *fiber.Ctx, err error) error {

	var code int
	var message string

	switch err.(type) {
	case *fiber.Error:
		code = err.(*fiber.Error).Code
		message = err.(*fiber.Error).Message
		break
	case *protocol_error.Error:
		message, code = protocol_error.Handle(err.(*protocol_error.Error))
	default:
		code = 400
		message = "Undefined Error"
		break
	}

	ctx.Status(code)
	return ctx.JSON(map[string]string{
		"error": message,
	})
}
