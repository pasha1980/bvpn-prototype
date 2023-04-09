package http

import (
	"bvpn-prototype/internal/protocols/entity"
	"github.com/gofiber/fiber/v2"
)

func Node(ctx *fiber.Ctx) *entity.Node {
	ip, ok := ctx.GetReqHeaders()["X-Forwarded-For"]
	if !ok {
		ip = ctx.IP()
	}

	return &entity.Node{
		IP: ip,
	}
}
