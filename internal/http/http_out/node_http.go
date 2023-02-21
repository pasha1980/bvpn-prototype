package http_out

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
)

func NewMempoolRecord(stored block_data.ChainStored, nodes ...entity.Node) error {
	return nil
}

func GetFullChain(node entity.Node) ([]entity.Block, error) {
	return nil, nil
}
