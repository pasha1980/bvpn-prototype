package api_in

import (
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/vpn/api_in/dto"
	"bvpn-prototype/internal/vpn/domain"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type VpnController struct {
}

func (VpnController) CreateConnection(ctx *fiber.Ctx) error {
	vpnService := di.Get("vpn_service").(domain.VpnService)
	profile, err := vpnService.CreateConnection()
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(dto.PublicProfileToDTO(*profile))
}

func (VpnController) BreakConnection(ctx *fiber.Ctx) error {
	// todo:
	return nil
}

func NewVpnController() VpnController {
	return VpnController{}
}
