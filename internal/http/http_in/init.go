package http_in

import (
	"github.com/gofiber/fiber/v2"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func InitHttp(addr string, tlsConfig *TLSConfig) error {
	app := fiber.New(fiber.Config{
		AppName:               "BVPN Prototype",
		ErrorHandler:          errorHandler,
		DisableStartupMessage: true,
	})

	app.Post("/:method", entrypoint)
	app.Get("/:method", entrypoint)
	app.Put("/:method", entrypoint)
	app.Patch("/:method", entrypoint)

	var err error
	if tlsConfig == nil {
		err = app.Listen(addr)
	} else {
		err = app.ListenTLS(addr, tlsConfig.CertFile, tlsConfig.KeyFile)
	}

	return err
}
