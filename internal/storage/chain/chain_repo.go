package chain

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/signer"
	"bvpn-prototype/internal/storage/chain/models"
	"bvpn-prototype/internal/storage/config"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"sort"
)

type ChainRepo struct {
	db *gorm.DB
}

func (r *ChainRepo) getData(block *entity.Block) error {
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

	var statuses []models.NodeStatus
	err = r.db.Where(&models.NodeStatus{
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

func (r *ChainRepo) saveData(block *entity.Block) error {
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
		case block_data.TypeNodeStatus:
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
			return errors.New("Invalid data type")
		}
	}

	return nil
}

func (r *ChainRepo) GetLastBlock() (*entity.Block, error) {
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

func (r *ChainRepo) GetBlockByHash(hash string) (*entity.Block, error) {
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

func (r *ChainRepo) GetBlockByNumber(number uint64) (*entity.Block, error) {
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

func (r *ChainRepo) GetChain(limit int, offset int) ([]entity.Block, error) {
	var blockModels []models.Block
	err := r.db.Limit(limit).Offset(offset).Find(&blockModels).Error
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

func (r *ChainRepo) SaveBlock(block entity.Block) (*entity.Block, error) {
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

func (r *ChainRepo) ReplaceChain(chain []entity.Block) error {
	tx := r.db.Raw("truncate table `chain`;")

	tx.Raw("truncate table `tx`;")
	tx.Raw("truncate table `offer`;")
	tx.Raw("truncate table `node_status`;")
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
			case block_data.TypeNodeStatus:
				tx.Save(models.NodeStatusToModel(data, uint(block.Number)))
				break
			case block_data.TypeTraffic:
				tx.Save(models.TrafficToModel(data, uint(block.Number)))
				break
			default:
				return errors.New("Invalid data type")
			}
		}
	}

	return tx.Error
}

func (r *ChainRepo) GetMyUTXOs() ([]block_data.ChainStored, error) {
	var utxos []block_data.ChainStored
	var utxoModels []models.Transaction

	myAddr := signer.GetAddr()
	sub := r.db.Select("o. from").Table("tx").Where("o.to = ?", myAddr)
	err := r.db.Where(models.Transaction{To: myAddr}).Where("ref not in (?)", sub).Find(&utxoModels).Error
	if err != nil {
		return nil, err
	}

	for _, model := range utxoModels {
		utxos = append(utxos, model.ToChainStored())
	}

	return utxos, nil
}

func NewChainRepo() *ChainRepo {
	db, err := gorm.Open(sqlite.Open(config.Get().StorageDirectory+"/chain.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Block{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.Traffic{})
	db.AutoMigrate(&models.NodeStatus{})
	db.AutoMigrate(&models.Offer{})

	return &ChainRepo{
		db: db,
	}
}
