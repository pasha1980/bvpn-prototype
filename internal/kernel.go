package internal

import (
	"bvpn-prototype/internal/protocols/entity"
)

type Kernel struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte

	PublicIp  string
	ChainPort uint

	Nodes []entity.Node
}

func (k *Kernel) Run() {

}
