package chain_http

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

const chainApiName = "BVPN Prototype Chain API"

type ChainNetworkManager struct {
	Port           string
	Client         http.Client
	ConnectedPeers []Peer
}

func (a *ChainNetworkManager) InitServer() {
	api := fiber.New(fiber.Config{
		AppName:               chainApiName,
		DisableStartupMessage: true,
		ErrorHandler:          a.errorHandler,
		ReduceMemoryUsage:     true,
	})

	api.Post("", a.entryPoint)

	err := api.Listen(a.Port)
	if err != nil {
		log.Fatal("Error while starting chain API server: " + err.Error())
	}
}

func (a *ChainNetworkManager) entryPoint(ctx *fiber.Ctx) error {
	var request apiRequest
	err := json.Unmarshal(ctx.Body(), &request)
	if err != nil {
		return err
	}

	function := apiFunctions[request.Method]
	if function == nil {
		return errors.New("Method `" + request.Method + "` not exist")
	}

	data, err := function(request.Arguments)
	if err != nil {
		return err
	}

	return ctx.JSON(
		apiResponse{
			Data: data,
		},
	)
}

func (a *ChainNetworkManager) errorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Status(http.StatusBadRequest)
	message := err.Error()
	return ctx.JSON(
		apiResponse{
			Error: &message,
		},
	)
}
