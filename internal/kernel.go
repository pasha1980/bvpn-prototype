package internal

import (
	"bvpn-prototype/internal/chain/chain_http"
)

type Kernel struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte

	WalletAddress string

	PublicIp  string
	ChainPort uint

	Peers []chain_http.Peer
}

func (k *Kernel) Run() {

}
