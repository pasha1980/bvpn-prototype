package permatent_tasks

import (
	"bvpn-prototype/internal/protocols"
	"time"
)

func Init() {
	go updateChain()
	go validateChain()
	go checkPeers()
}

func updateChain() {
	for range time.Tick(10 * time.Minute) {
		protocols.GetChainProtocol().UpdateChain()
	}
}

func validateChain() {
	for range time.Tick(time.Minute) {
		protocols.GetChainProtocol().ValidateChain() // todo: what to do
	}
}

func checkPeers() {
	for range time.Tick(time.Hour) {
		protocols.GetPeerProtocol().CheckPeers()
	}
}
