package http

import (
	"github.com/gofiber/fiber/v2"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func Init(port string, tlsConfig *TLSConfig) error {
	app := fiber.New(fiber.Config{
		AppName:               "BVPN Prototype",
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
	})

	fillChainEntrypoints(app)
	fillVpnEntrypoints(app)
	fillPeerEntrypoints(app)

	var err error
	if tlsConfig == nil {
		err = app.Listen(":" + port)
	} else {
		err = app.ListenTLS(":"+port, tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return err
}

func fillChainEntrypoints(app *fiber.App) {
	app.Post("/chain/:method", chainEntrypoint)
	app.Get("/chain/:method", chainEntrypoint)
	app.Put("/chain/:method", chainEntrypoint)
	app.Patch("/chain/:method", chainEntrypoint)
}

func fillVpnEntrypoints(app *fiber.App) {
	app.Post("/vpn/:method", vpnEntrypoint)
	app.Get("/vpn/:method", vpnEntrypoint)
	app.Put("/vpn/:method", vpnEntrypoint)
	app.Patch("/vpn/:method", vpnEntrypoint)
}

func fillPeerEntrypoints(app *fiber.App) {
	app.Post("/peer/:method", peerEntrypoint)
	app.Get("/peer/:method", peerEntrypoint)
	app.Put("/peer/:method", peerEntrypoint)
	app.Patch("/peer/:method", peerEntrypoint)
}
