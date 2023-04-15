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

	initChainHttp(app)
	initVpnHttp(app)
	initPeerHttp(app)

	var err error
	if tlsConfig == nil {
		err = app.Listen(port)
	} else {
		err = app.ListenTLS(port, tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return err
}

func initChainHttp(app *fiber.App) {
	app.Post("/chain/:method", chainEntrypoint)
	app.Get("/chain/:method", chainEntrypoint)
	app.Put("/chain/:method", chainEntrypoint)
	app.Patch("/chain/:method", chainEntrypoint)
}

func initVpnHttp(app *fiber.App) {
	app.Post("/vpn/:method", vpnEntrypoint)
	app.Get("/vpn/:method", vpnEntrypoint)
	app.Put("/vpn/:method", vpnEntrypoint)
	app.Patch("/vpn/:method", vpnEntrypoint)
}

func initPeerHttp(app *fiber.App) {
	app.Post("/vpn/:method", peerEntrypoint)
	app.Get("/vpn/:method", peerEntrypoint)
	app.Put("/vpn/:method", peerEntrypoint)
	app.Patch("/vpn/:method", peerEntrypoint)
}
