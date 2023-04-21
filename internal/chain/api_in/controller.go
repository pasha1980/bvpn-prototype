package api_in

import (
	"bvpn-prototype/internal/chain/api_in/dto"
	"bvpn-prototype/internal/chain/domain"
	"bvpn-prototype/internal/common"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/errors/http_errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ChainController struct {
}

func (ChainController) AddTx(ctx *fiber.Ctx) error {
	var tx dto.Transaction

	err := ctx.BodyParser(&tx)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(tx); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	err = c.AddToMempool(tx.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func (ChainController) AddOffer(ctx *fiber.Ctx) error {
	var offer dto.Offer

	err := ctx.BodyParser(&offer)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(offer); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	err = c.AddToMempool(offer.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func (ChainController) AddTraffic(ctx *fiber.Ctx) error {
	var traffic dto.Traffic

	err := ctx.BodyParser(&traffic)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(traffic); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	err = c.AddToMempool(traffic.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func (ChainController) AddConnectionBreak(ctx *fiber.Ctx) error {
	var connectionBreak dto.ConnectionBreak

	err := ctx.BodyParser(&connectionBreak)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(connectionBreak); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	err = c.AddToMempool(connectionBreak.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func (ChainController) AddBlock(ctx *fiber.Ctx) error {
	var blockDto dto.BlockDto

	err := ctx.BodyParser(&blockDto)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(blockDto); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	err = c.AddBlock(blockDto.ToEntity(), Node(ctx))
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func (ChainController) GetChain(ctx *fiber.Ctx) error {
	var body dto.PaginationDto
	err := ctx.BodyParser(&body)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(body); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	c := di.Get("chain_service").(domain.ChainService)
	chain, err := c.GetChain(body.Limit, body.Offset)
	if err != nil {
		return err
	}

	var blockDtos []dto.BlockDto
	for _, block := range chain {
		blockDtos = append(blockDtos, dto.BlockToDto(block))
	}

	chainDto := dto.ChainDto{
		Chain:      blockDtos,
		TotalCount: len(blockDtos),
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(chainDto))
}

func NewChainController() ChainController {
	return ChainController{}
}
