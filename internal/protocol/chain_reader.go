package protocol

import "bvpn-prototype/internal/protocol/entity"

type ChainReader interface {
	Start()
	Next() *entity.Block
	Last() *entity.Block
	Len() int64
}
