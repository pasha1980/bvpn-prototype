package chain

import (
	"bvpn-prototype/internal/protocols/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type ChainRepo struct {
	db *gorm.DB
}

func (r *ChainRepo) GetLastBlock() *entity.Block {
	var blockModel BlockModel
	r.db.Last(&blockModel)
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetBlockByHash(hash string) *entity.Block {
	var blockModel BlockModel
	r.db.Where(&BlockModel{
		hash: hash,
	}).Find(&blockModel)
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetBlockByNumber(number uint64) *entity.Block {
	var blockModel BlockModel
	r.db.Find(&blockModel, uint(number))
	return blockModel.modelToEntity()
}

func (r *ChainRepo) GetFullChain() []entity.Block {
	var blockModels []BlockModel
	r.db.Find(&blockModels)

	var blocks []entity.Block
	for _, model := range blockModels {
		blocks = append(blocks, *model.modelToEntity())
	}

	return blocks
}

func (r *ChainRepo) SaveBlock(block *entity.Block) *entity.Block {
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
