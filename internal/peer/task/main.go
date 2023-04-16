package task

import (
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/peer/domain"
	"time"
)

func Run() {
	go checkPeers()
}

func checkPeers() {
	peerService := di.Get("peer_service").(domain.PeerService)
	for range time.Tick(time.Hour) {
		peerService.CheckPeers()
	}
}
