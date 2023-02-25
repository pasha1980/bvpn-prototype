package http_in

import (
	"bvpn-prototype/internal/logger"
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
		switch err.(*protocol_error.Error).Code {
		case protocol_error.MessageErrorCode:
			code = 400
			message = err.Error()
			break
		case protocol_error.LogErrorCode:
			code = 400
			message = err.Error()
			logger.LogError(message)
			break
		case protocol_error.LogInternalErrorCode:
			code = 500
			message = err.Error()
			logger.LogError(message)
			break
		}
		break
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
