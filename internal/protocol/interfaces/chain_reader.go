package interfaces

import "bvpn-prototype/internal/protocol/entity"

type ChainReader interface {
	Start()
	Next() *entity.Block
	Last() *entity.Block
	Previous(number uint64) *entity.Block
	Len() int64
}
