package api_in

import (
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/http/http_errors"
	"bvpn-prototype/internal/peer/api_in/dto"
	"bvpn-prototype/internal/peer/domain"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type PeerController struct {
}

func (c PeerController) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func (c PeerController) AddPeer(ctx *fiber.Ctx) error {
	var body dto.PeerDto
	err := ctx.BodyParser(&body)
	if err != nil {
		return http_errors.ErrorInvalidRequest
	}

	peerService := di.Get("peer_service").(domain.PeerService)
	err = peerService.AddPeer(body.ToEntity())
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(map[string]bool{
		"success": true,
	})
}

func NewPeerController() PeerController {
	return PeerController{}
}
