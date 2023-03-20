package repo

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
)

type ChainStorageRepo interface {
	GetLastBlock() (*entity.Block, error)
	GetBlockByHash(hash string) (*entity.Block, error)
	GetBlockByNumber(number uint64) (*entity.Block, error)
	SaveBlock(block entity.Block) (*entity.Block, error)
	GetChain(limit int, offset int) ([]entity.Block, error)

	ReplaceChain(chain []entity.Block) error

	GetMyUTXOs() ([]block_data.ChainStored, error)
}
