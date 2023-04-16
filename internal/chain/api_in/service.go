package api_in

import (
	"bvpn-prototype/internal/protocol/entity"
	"github.com/gofiber/fiber/v2"
)

func Node(ctx *fiber.Ctx) *entity.Node {
	return &entity.Node{
		IP: ctx.IP(),
	}
}
