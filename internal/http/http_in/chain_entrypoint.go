package http_in

import (
	"bvpn-prototype/internal/api"
	"bvpn-prototype/internal/http/http_errors"
	"errors"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

func entrypoint(ctx *fiber.Ctx) error {
	c := api.HTTP()
	methodParam := ctx.Params("method")
	firstLetter := strings.ToUpper(methodParam[0:1])
	methodName := firstLetter + methodParam[1:]
	method := reflect.ValueOf(*c).MethodByName(methodName)

	if !method.IsValid() {
		return http_errors.ErrorMethodNotFound
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(ctx)

	response := method.Call(in)[0]
	if response.IsZero() || response.IsValid() {
		return nil
	}

	return errors.New(response.MethodByName("Error").Call([]reflect.Value{})[0].String())
}
