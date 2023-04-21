package api_in

import (
	"bvpn-prototype/internal/common"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/errors/http_errors"
	"bvpn-prototype/internal/vpn/api_in/dto"
	"bvpn-prototype/internal/vpn/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
)

type VpnController struct {
}

func (VpnController) CreateConnection(ctx *fiber.Ctx) error {
	var d dto.InitDTO

	err := ctx.BodyParser(&d)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(d); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	vpnService := di.Get("vpn_service").(domain.VpnService)
	profile, err := vpnService.CreateConnection(d.Addr)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(common.MakeHttpPage(dto.PublicProfileToDTO(*profile)))
}

func (VpnController) BreakConnection(ctx *fiber.Ctx) error {
	var d dto.ProfileShortDto

	err := ctx.BodyParser(&d)
	if err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	if err = validator.New().Struct(d); err != nil {
		return http_errors.InvalidRequest(err.Error())
	}

	id, _ := uuid.Parse(d.ID)
	vpnService := di.Get("vpn_service").(domain.VpnService)
	err = vpnService.BreakConnection(id)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(common.MakeHttpPage(true))
}

func NewVpnController() VpnController {
	return VpnController{}
}
