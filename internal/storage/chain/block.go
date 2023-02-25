package chain

import (
	"bvpn-prototype/internal/protocols/entity"
	"encoding/json"
	"time"
)

type BlockModel struct {
	ID           uint `gorm:"PRIMARY_KEY"`
	Hash         string
	PreviousHash string
	Data         []byte
	Timestamp    time.Time
}

func (b BlockModel) TableName() string {
	return "chain"
}

func (b *BlockModel) modelToEntity() *entity.Block {
	e := entity.Block{
		Number:       uint64(b.ID),
		Hash:         b.Hash,
		PreviousHash: b.PreviousHash,
		TimeStamp:    b.Timestamp,
	}

	json.Unmarshal(b.Data, &e.Data)
	return &e
}

func blockToModel(block entity.Block) *BlockModel {
	model := BlockModel{
		Hash:         block.Hash,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.TimeStamp,
	}

	data, _ := json.Marshal(block.Data)
	model.Data = data

	if block.Number != 0 {
		model.ID = uint(block.Number)
	}

	return &model
}
