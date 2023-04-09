package http

import (
	"bvpn-prototype/internal/http/http_dto"
	"bvpn-prototype/internal/http/http_dto/mempool_data_dto"
	"bvpn-prototype/internal/http/http_errors"
	"bvpn-prototype/internal/protocols"
	"bvpn-prototype/internal/protocols/signer"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type HttpApi struct {
	ChainProtocol *protocols.ChainProtocol
	PeerProtocol  *protocols.PeerProtocol
}

func (HttpApi) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) AddTx(ctx *fiber.Ctx) error {
	var tx mempool_data_dto.Transaction

	err := ctx.BodyParser(&tx)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(tx.ToEntity(), Node(ctx))

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) AddOffer(ctx *fiber.Ctx) error {
	var offer mempool_data_dto.Offer

	err := ctx.BodyParser(&offer)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(offer.ToEntity(), Node(ctx))

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) AddTraffic(ctx *fiber.Ctx) error {
	var traffic mempool_data_dto.Traffic

	err := ctx.BodyParser(&traffic)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(traffic.ToEntity(), Node(ctx))

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) AddConnectionBreak(ctx *fiber.Ctx) error {
	var connectionBreak mempool_data_dto.ConnectionBreak

	err := ctx.BodyParser(&connectionBreak)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	c.ChainProtocol.AddToMempool(connectionBreak.ToEntity(), Node(ctx))

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) AddBlock(ctx *fiber.Ctx) error {
	var blockDto http_dto.BlockDto

	err := ctx.BodyParser(&blockDto)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	err = c.ChainProtocol.AddBlock(blockDto.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c HttpApi) GetChain(ctx *fiber.Ctx) error {
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

func (c HttpApi) AddPeer(ctx *fiber.Ctx) error {
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

func (c HttpApi) GetAddress(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(map[string]string{
		"address": signer.GetAddr(),
	})
}
