package http_in

import (
	"bvpn-prototype/internal/http/http_dto"
	"bvpn-prototype/internal/http/http_dto/mempool_data_dto"
	"bvpn-prototype/internal/http/http_errors"
	"bvpn-prototype/internal/protocols"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"reflect"
)

type HttpController struct {
	ChainProtocol protocols.ChainProtocol
	PeerProtocol  protocols.PeerProtocol
}

func (c *HttpController) healthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) addTx(ctx *fiber.Ctx) error {
	var tx mempool_data_dto.Transaction

	err := ctx.BodyParser(&tx)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(tx.ToEntity())

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
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

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
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

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
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

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) addBlock(ctx *fiber.Ctx) error {
	var blockDto http_dto.BlockDto

	err := ctx.BodyParser(&blockDto)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	err = c.ChainProtocol.AddBlock(blockDto.ToEntity())
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c *HttpController) getChain(ctx *fiber.Ctx) error {
	var body map[string]int
	err := ctx.BodyParser(&body)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	chain, err := c.ChainProtocol.GetChain(body["limit"], body["offset"]) // todo: unsafe
	if err != nil {
		return err
	}

	var blockDtos []http_dto.BlockDto
	for _, block := range chain {
		blockDtos = append(blockDtos, http_dto.BlockToDto(block))
	}

	chainDto := http_dto.ChainDto{
		Chain:      blockDtos,
		TotalCount: len(blockDtos),
	}

	return ctx.Status(http.StatusOK).JSON(chainDto)
}

func (c *HttpController) addPeer(ctx *fiber.Ctx) error {
	var body http_dto.PeerDto
	err := ctx.BodyParser(&body)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	err = c.PeerProtocol.AddNewPeer(body.ToEntity())
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
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
