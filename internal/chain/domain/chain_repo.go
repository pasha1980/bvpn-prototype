package domain

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
)

type ChainRepository interface {
	GetLastBlock() (*entity.Block, error)
	GetBlockByHash(hash string) (*entity.Block, error)
	GetBlockByNumber(number uint64) (*entity.Block, error)
	GetChain(limit *int, offset *int) ([]entity.Block, error)
	SaveBlock(block entity.Block) (*entity.Block, error)
	ReplaceChain(chain []entity.Block) error
	GetUTXOs(addr string) ([]block_data.ChainStored, error)
	GetLastOffer(pubKey string) (*block_data.ChainStored, error)
}
