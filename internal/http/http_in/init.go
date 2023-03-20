package http_in

import (
	"github.com/gofiber/fiber/v2"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func InitHttp(controller HttpController, addr string, tlsConfig *TLSConfig) error {
	app := fiber.New(fiber.Config{
		AppName:               "BVPN Prototype",
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
	})

	app.Post("/:method", controller.HttpEntrypoint)
	app.Get("/:method", controller.HttpEntrypoint)
	app.Put("/:method", controller.HttpEntrypoint)
	app.Patch("/:method", controller.HttpEntrypoint)

	var err error
	if tlsConfig == nil {
		err = app.Listen(addr)
	} else {
		err = app.ListenTLS(addr, tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return err
}
