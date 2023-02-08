package chain_storage

import (
	"bvpn-prototype/internal/chain/chain_domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type ChainRepo struct {
	db *gorm.DB
}

func (r *ChainRepo) GetLastBlock() *chain_domain.Block {
	var blockModel BlockModel
	r.db.Last(&blockModel)
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetBlockByHash(hash string) *chain_domain.Block {
	var blockModel BlockModel
	r.db.Where(&BlockModel{
		hash: hash,
	}).Find(&blockModel)
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetBlockByNumber(number uint64) *chain_domain.Block {
	var blockModel BlockModel
	r.db.Find(&blockModel, uint(number))
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetFullChain() []chain_domain.Block {
	var blockModels []BlockModel
	r.db.Find(&blockModels)

	var blocks []chain_domain.Block
	for _, model := range blockModels {
		blocks = append(blocks, *model.modelToEntity())
	}

	return blocks
}

func (r *ChainRepo) SaveBlock(block *chain_domain.Block) *chain_domain.Block {
	model := blockToModel(*block)
	r.db.Save(model)
	return model.modelToEntity()
}

func NewChainRepo() *ChainRepo {
	db, err := gorm.Open(sqlite.Open("/opt/bvpn/chain.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	return &ChainRepo{
		db: db,
	}
}