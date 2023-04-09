package api

import (
	"bvpn-prototype/internal/api/cli"
	"bvpn-prototype/internal/api/http"
	"bvpn-prototype/internal/protocols"
	"bvpn-prototype/internal/storage/config"
)

func CLI(cfg *config.Config) *cli.CliApi {
	return &cli.CliApi{
		ChainProtocol: protocols.GetChainProtocol(),
		PeerProtocol:  protocols.GetPeerProtocol(),
		Config:        cfg,
	}
}

func HTTP() *http.HttpApi {
	return &http.HttpApi{
		ChainProtocol: protocols.GetChainProtocol(),
		PeerProtocol:  protocols.GetPeerProtocol(),
	}
}
