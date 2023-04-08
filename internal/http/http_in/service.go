package http_in

import (
	"bvpn-prototype/internal/protocols/entity"
	"github.com/gofiber/fiber/v2"
)

func Node(ctx *fiber.Ctx) *entity.Node {
	return &entity.Node{
		IP: ctx.IP(),
	}
}
