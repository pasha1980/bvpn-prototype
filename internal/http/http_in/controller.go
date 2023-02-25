package http_in

import (
	"bvpn-prototype/internal/http/http_dto/mempool_data_dto"
	"bvpn-prototype/internal/http/http_errors"
	"bvpn-prototype/internal/protocols"
	"errors"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type HttpController struct {
	ChainProtocol protocols.ChainProtocol
}

func (c *HttpController) addTx(ctx *fiber.Ctx) error {
	var tx mempool_data_dto.Transaction

	err := ctx.BodyParser(&tx)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(tx.ToEntity())

	return ctx.JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) addOffer(ctx *fiber.Ctx) error {
	var offer mempool_data_dto.Offer

	err := ctx.BodyParser(&offer)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(offer.ToEntity())

	return ctx.JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) addTraffic(ctx *fiber.Ctx) error {
	var traffic mempool_data_dto.Traffic

	err := ctx.BodyParser(&traffic)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(traffic.ToEntity())

	return ctx.JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) addNodeStatus(ctx *fiber.Ctx) error {
	var nodeStatus mempool_data_dto.NodeStatus

	err := ctx.BodyParser(&nodeStatus)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(nodeStatus.ToEntity())

	return ctx.JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) HttpEntrypoint(ctx *fiber.Ctx) error {
	methodName := ctx.Params("method")
	method := reflect.ValueOf(c).MethodByName(methodName)
	if method.IsZero() || method.IsValid() {
		return http_errors.ErrorMethodNotFound
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(ctx)

	response := method.Call(in)[0]
	if response.IsZero() || response.IsValid() {
		return nil
	}

	return errors.New(response.MethodByName("Error").Call([]reflect.Value{})[0].String()) // todo: test
}
