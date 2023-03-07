package chain

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type ChainRepo struct {
	db *gorm.DB
}

func (r *ChainRepo) GetLastBlock() (*entity.Block, error) {
	var blockModel BlockModel
	err := r.db.Last(&blockModel).Error
	if err != nil {
		return nil, err
	}

	return blockModel.modelToEntity(), nil
}

func (r *ChainRepo) GetBlockByHash(hash string) (*entity.Block, error) {
	var blockModel BlockModel
	err := r.db.Where(&BlockModel{
		Hash: hash,
	}).Find(&blockModel).Error
	if err != nil {
		return nil, err
	}

	return blockModel.modelToEntity(), nil
}

func (r *ChainRepo) GetBlockByNumber(number uint64) (*entity.Block, error) {
	var blockModel BlockModel
	err := r.db.Find(&blockModel, uint(number)).Error
	if err != nil {
		return nil, err
	}

	return blockModel.modelToEntity(), nil
}

func (r *ChainRepo) GetChain(limit int, offset int) ([]entity.Block, error) {
	var blockModels []BlockModel
	err := r.db.Limit(limit).Offset(offset).Find(&blockModels).Error
	if err != nil {
		return nil, err
	}

	var blocks []entity.Block
	for _, model := range blockModels {
		blocks = append(blocks, *model.modelToEntity())
	}

	return blocks, nil
}

func (r *ChainRepo) SaveBlock(block entity.Block) (*entity.Block, error) {
	model := blockToModel(block)
	err := r.db.Save(model).Error
	if err != nil {
		return nil, err
	}

	return model.modelToEntity(), err
}

func (r *ChainRepo) ReplaceChain(chain []entity.Block) error {
	tx := r.db.Raw("truncate table `chain`")

	for _, block := range chain {
		tx.Save(blockToModel(block))
	}

	return tx.Error
}

func NewChainRepo() *ChainRepo {
	db, err := gorm.Open(sqlite.Open(config.Get().StorageDirectory+"/chain.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	return &ChainRepo{
		db: db,
	}
}
