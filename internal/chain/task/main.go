package task

import (
	"bvpn-prototype/internal/chain/domain"
	"bvpn-prototype/internal/infrastructure/di"
	"time"
)

func Run() {
	go updateChain()
	go validateChain()
}

func updateChain() {
	chainService := di.Get("chain_service").(domain.ChainService)
	for range time.Tick(10 * time.Minute) {
		chainService.UpdateChain()
	}
}

func validateChain() {
	chainService := di.Get("chain_service").(domain.ChainService)
	for range time.Tick(time.Minute) {
		chainService.ValidateStoredChain()
	}
}
