package mempool

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
)

type pool struct {
	Data map[string]block_data.ChainStored
}

var storage *pool
