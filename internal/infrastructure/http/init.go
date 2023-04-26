package http

import (
	"bvpn-prototype/internal/common"
	"bvpn-prototype/internal/infrastructure/config"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func Init(port string, tlsConfig *TLSConfig) error {
	app := fiber.New(fiber.Config{
		AppName:               common.PROJECT_NAME,
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
	})

	fillChainEntrypoint(app)
	fillVpnEntrypoint(app)
	fillPeerEntrypoint(app)

	version(app)

	var err error
	if tlsConfig == nil {
		err = app.Listen(":" + port)
	} else {
		err = app.ListenTLS(":"+port, tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return err
}

func fillChainEntrypoint(app *fiber.App) {
	app.Post("/chain/:method", chainEntrypoint)
	app.Get("/chain/:method", chainEntrypoint)
	app.Put("/chain/:method", chainEntrypoint)
	app.Patch("/chain/:method", chainEntrypoint)
}

func fillVpnEntrypoint(app *fiber.App) {
	app.Post("/vpn/:method", vpnEntrypoint)
	app.Get("/vpn/:method", vpnEntrypoint)
	app.Put("/vpn/:method", vpnEntrypoint)
	app.Patch("/vpn/:method", vpnEntrypoint)
}

func fillPeerEntrypoint(app *fiber.App) {
	app.Post("/peer/:method", peerEntrypoint)
	app.Get("/peer/:method", peerEntrypoint)
	app.Put("/peer/:method", peerEntrypoint)
	app.Patch("/peer/:method", peerEntrypoint)
}

func version(app *fiber.App) {
	function := func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).Format(config.VERSION)
	}

	app.Post("/version", function)
	app.Get("/version", function)
	app.Put("/version", function)
	app.Patch("/version", function)
}
