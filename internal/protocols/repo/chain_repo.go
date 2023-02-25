package repo

import (
	"bvpn-prototype/internal/protocols/entity"
)

type ChainStorageRepo interface {
	GetLastBlock() (*entity.Block, error)
	GetBlockByHash(hash string) (*entity.Block, error)
	GetBlockByNumber(number uint64) (*entity.Block, error)
	SaveBlock(block entity.Block) (*entity.Block, error)

	ReplaceChain(chain []entity.Block) error
}
