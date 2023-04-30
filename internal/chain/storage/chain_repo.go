package storage

import (
	"bvpn-prototype/internal/chain/storage/models"
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sort"
)

type ChainRepository struct {
	db *gorm.DB

	readerIndex uint64
}

func (r *ChainRepository) getData(block *entity.Block) error {
	var data []block_data.ChainStored

	var txs []models.Transaction
	err := r.db.Where(&models.Transaction{
		BlockID: uint(block.Number),
	}).Find(&txs).Error
	if err != nil {
		return err
	}
	for _, tx := range txs {
		data = append(data, tx.ToChainStored())
	}

	var offers []models.Offer
	err = r.db.Where(&models.Offer{
		BlockID: uint(block.Number),
	}).Find(&offers).Error
	if err != nil {
		return err
	}
	for _, offer := range offers {
		data = append(data, offer.ToChainStored())
	}

	var statuses []models.ConnectionBreak
	err = r.db.Where(&models.ConnectionBreak{
		BlockID: uint(block.Number),
	}).Find(&statuses).Error
	if err != nil {
		return err
	}
	for _, status := range statuses {
		data = append(data, status.ToChainStored())
	}

	var traffic []models.Traffic
	err = r.db.Where(&models.Traffic{
		BlockID: uint(block.Number),
	}).Find(&traffic).Error
	if err != nil {
		return err
	}
	for _, tr := range traffic {
		data = append(data, tr.ToChainStored())
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].ID.ID() < data[j].ID.ID()
	})

	block.Data = data
	return nil
}

func (r *ChainRepository) saveData(block *entity.Block) error {
	for _, data := range block.Data {
		switch data.Type {
		case block_data.TypeTransaction:
			model := models.TxToModel(data, uint(block.Number))
			err := r.db.Save(model).Error
			if err != nil {
				return err
			}
			break
		case block_data.TypeOffer:
			model := models.OfferToModel(data, uint(block.Number))
			err := r.db.Save(model).Error
			if err != nil {
				return err
			}
			break
		case block_data.TypeConnectionBreak:
			model := models.NodeStatusToModel(data, uint(block.Number))
			err := r.db.Save(model).Error
			if err != nil {
				return err
			}
			break
		case block_data.TypeTraffic:
			model := models.TrafficToModel(data, uint(block.Number))
			err := r.db.Save(model).Error
			if err != nil {
				return err
			}
			break
		default:
			model := models.UndefinedDataModel(data, uint(block.Number))
			err := r.db.Save(model).Error
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (r *ChainRepository) GetLastBlock() (*entity.Block, error) {
	var blockModel models.Block
	err := r.db.Last(&blockModel).Error
	if err != nil {
		return nil, err
	}

	b := blockModel.ModelToEntity()
	err = r.getData(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *ChainRepository) GetBlockByHash(hash string) (*entity.Block, error) {
	var blockModel models.Block
	err := r.db.Where(&models.Block{
		Hash: hash,
	}).Find(&blockModel).Error
	if err != nil {
		return nil, err
	}

	b := blockModel.ModelToEntity()
	err = r.getData(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *ChainRepository) GetBlockByNumber(number uint64) (*entity.Block, error) {
	var blockModel models.Block
	err := r.db.Find(&blockModel, uint(number)).Error
	if err != nil {
		return nil, err
	}

	b := blockModel.ModelToEntity()
	err = r.getData(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *ChainRepository) GetChain(limit *int, offset *int) ([]entity.Block, error) {
	var blockModels []models.Block
	query := r.db

	if limit != nil {
		query.Limit(*limit)
		if offset != nil {
			query.Offset(*offset)
		}
	}

	err := query.Find(&blockModels).Error
	if err != nil {
		return nil, err
	}

	var blocks []entity.Block
	for _, model := range blockModels {
		b := model.ModelToEntity()
		err = r.getData(b)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, *b)
	}

	return blocks, nil
}

func (r *ChainRepository) SaveBlock(block entity.Block) (*entity.Block, error) {
	model := models.BlockToModel(block)
	err := r.db.Save(model).Error
	if err != nil {
		return nil, err
	}

	err = r.saveData(&block)
	if err != nil {
		return nil, err
	}

	return model.ModelToEntity(), err
}

func (r *ChainRepository) ReplaceChain(chain []entity.Block) error {
	tx := r.db.Raw("truncate table `chain`;")

	tx.Raw("truncate table `tx`;")
	tx.Raw("truncate table `offer`;")
	tx.Raw("truncate table `connection_break`;")
	tx.Raw("truncate table `traffic`;")

	for _, block := range chain {
		tx.Save(models.BlockToModel(block))
		for _, data := range block.Data {
			switch data.Type {
			case block_data.TypeTransaction:
				tx.Save(models.TxToModel(data, uint(block.Number)))
				break
			case block_data.TypeOffer:
				tx.Save(models.OfferToModel(data, uint(block.Number)))
				break
			case block_data.TypeConnectionBreak:
				tx.Save(models.NodeStatusToModel(data, uint(block.Number)))
				break
			case block_data.TypeTraffic:
				tx.Save(models.TrafficToModel(data, uint(block.Number)))
				break
			default:
				tx.Save(models.UndefinedDataModel(data, uint(block.Number)))
			}
		}
	}

	return tx.Error
}

func (r *ChainRepository) GetUTXOs(addr string) ([]block_data.ChainStored, error) {
	var utxos []block_data.ChainStored
	var utxoModels []models.Transaction

	sub := r.db.Select("o. from").Table("tx").Where("o.to = ?", addr)
	err := r.db.Where(models.Transaction{To: addr}).Where("ref not in (?)", sub).Find(&utxoModels).Error
	if err != nil {
		return nil, err
	}

	for _, model := range utxoModels {
		utxos = append(utxos, model.ToChainStored())
	}

	return utxos, nil
}

func (r *ChainRepository) GetLastOffer(pubKey string) (*block_data.ChainStored, error) {
	var offerModel models.Offer

	err := r.db.Where(models.Offer{PubKey: pubKey}).Order("timestamp asc").Last(&offerModel).Error
	if err != nil {
		return nil, err
	}

	cs := offerModel.ToChainStored()
	return &cs, nil
}

func (r *ChainRepository) Start() {
	r.readerIndex = 1
}

func (r *ChainRepository) Next() *entity.Block {
	r.readerIndex++
	lastBlock, _ := r.GetLastBlock() // todo: errors
	r.getData(lastBlock)             // todo: errors
	if r.readerIndex > lastBlock.Number {
		return nil
	}

	block, _ := r.GetBlockByNumber(r.readerIndex) // todo: errors
	return block
}

func (r *ChainRepository) Last() *entity.Block {
	block, _ := r.GetLastBlock() // todo: errors
	r.getData(block)             // todo: errors
	return block
}

func (r *ChainRepository) Previous(number uint64) *entity.Block {
	block, _ := r.GetBlockByNumber(number - 1)
	r.getData(block) // todo: errors
	return block
}

func (r *ChainRepository) Len() int64 {
	var count int64
	r.db.Model(&models.Block{}).Count(&count)
	return count
}

func NewChainRepo() (*ChainRepository, error) {
	db, err := gorm.Open(sqlite.Open(config.Get().StorageDirectory+"/chain.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Block{},
		&models.Transaction{},
		&models.Traffic{},
		&models.ConnectionBreak{},
		&models.Offer{},
		&models.UndefinedData{},
	)

	if err != nil {
		return nil, err
	}

	return &ChainRepository{
		db:          db,
		readerIndex: 1,
	}, nil
}
