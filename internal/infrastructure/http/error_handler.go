package http

import (
	"bvpn-prototype/internal/infrastructure/common"
	"bvpn-prototype/internal/infrastructure/errors"
	"github.com/gofiber/fiber/v2"
)

func errorHandler(ctx *fiber.Ctx, err error) error {

	var statusCode int
	var message string
	var errorCode int
	var data any

	switch err.(type) {
	case *fiber.Error:
		statusCode = err.(*fiber.Error).Code
		errorCode = 10011
		message = err.(*fiber.Error).Message
		data = nil
		break
	case *errors.Error:
		statusCode = 400
		errStruct := err.(errors.Error)
		errorCode = errStruct.Code
		message = err.Error()
		data = errStruct.Data

		errStruct.Log()
		break
	default:
		statusCode = 500
		errorCode = 10011
		message = err.Error()
		data = nil
		break
	}

	ctx.Status(statusCode)
	return ctx.JSON(common.MakeErrorPage(errorCode, message, data))
}
