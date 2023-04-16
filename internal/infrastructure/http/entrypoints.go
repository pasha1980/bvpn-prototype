package http

import (
	chain_api "bvpn-prototype/internal/chain/api_in"
	"bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/infrastructure/errors/http_errors"
	peer_api "bvpn-prototype/internal/peer/api_in"
	vpn_api "bvpn-prototype/internal/vpn/api_in"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

func chainEntrypoint(ctx *fiber.Ctx) error {
	return rpcHandle(ctx, chain_api.NewChainController())
}

func vpnEntrypoint(ctx *fiber.Ctx) error {
	return rpcHandle(ctx, vpn_api.NewVpnController())
}

func peerEntrypoint(ctx *fiber.Ctx) error {
	return rpcHandle(ctx, peer_api.NewPeerController())
}

func rpcHandle(ctx *fiber.Ctx, controller interface{}) error {
	methodParam := ctx.Params("method")
	firstLetter := strings.ToUpper(methodParam[0:1])
	methodName := firstLetter + methodParam[1:]
	method := reflect.ValueOf(controller).MethodByName(methodName)

	if !method.IsValid() {
		return http_errors.MethodNotFoundHttpError("Method "+methodParam+" not found", nil)
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(ctx)

	response := method.Call(in)[0]
	if response.IsZero() || response.IsValid() {
		return nil
	}

	return response.Interface().(errors.Error)
}
