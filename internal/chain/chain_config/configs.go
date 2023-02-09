package chain_config

import "bvpn-prototype/internal/chain/chain_http"

type ChainConfig struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte

	WalletAddress string

	ConnectedPeers []chain_http.Peer
}

var Storage ChainConfig
