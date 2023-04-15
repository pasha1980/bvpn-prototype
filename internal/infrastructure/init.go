package infrastructure

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/infrastructure/http"
	"bvpn-prototype/internal/infrastructure/logger"
	"bvpn-prototype/internal/infrastructure/tasks"
)

const (
	WITH_HTTP  = 0x20
	WITH_TASKS = 0x30
)

func Init(c config.Config, flags int) error {
	config.Set(c)
	logger.Init()

	if flags&WITH_HTTP != 0 {
		err := http.Init(":"+c.HttpPort, nil) // todo: tls
		if err != nil {
			return err
		}
	}

	if flags&WITH_TASKS != 0 {
		tasks.Init()
	}

	return nil
}
