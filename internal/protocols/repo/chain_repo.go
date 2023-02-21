package repo

import (
	"bvpn-prototype/internal/protocols/entity"
)

type ChainStorageRepo interface {
	GetLastBlock() *entity.Block
	GetBlockByHash(hash string) *entity.Block
	GetBlockByNumber(number uint64) *entity.Block
	SaveBlock(block *entity.Block) *entity.Block

	ReplaceChain(chain []entity.Block) error
}
