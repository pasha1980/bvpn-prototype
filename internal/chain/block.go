package chain

import (
	"bvpn-prototype/internal/protocols/entity"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type BlockModel struct {
	gorm.Model
	hash         string
	previousHash string
	data         []byte
	timestamp    time.Time
}

func (b *BlockModel) modelToEntity() *entity.Block {
	entity := entity.Block{
		Number:       uint64(b.ID),
		Hash:         b.hash,
		PreviousHash: b.previousHash,
		TimeStamp:    b.timestamp,
	}

	json.Unmarshal(b.data, &entity.Data)
	return &entity
}

func blockToModel(block entity.Block) *BlockModel {
	model := BlockModel{
		hash:         block.Hash,
		previousHash: block.PreviousHash,
		timestamp:    block.TimeStamp,
	}

	data, _ := json.Marshal(block.Data)
	model.data = data

	if block.Number != 0 {
		model.ID = uint(block.Number)
	}

	return &model
}
