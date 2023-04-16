package models

import (
	"bvpn-prototype/internal/protocol/entity"
	"time"
)

type Block struct {
	ID           uint   `gorm:"PRIMARY_KEY"`
	Hash         string `gorm:"index"`
	PreviousHash string `gorm:"index"`
	Timestamp    time.Time
	CreatedBy    string `gorm:"index"`
	Next         string
}

func (b Block) TableName() string {
	return "chain"
}

func (b *Block) ModelToEntity() *entity.Block {
	e := entity.Block{
		Number:       uint64(b.ID),
		Hash:         b.Hash,
		PreviousHash: b.PreviousHash,
		TimeStamp:    b.Timestamp,
		Next:         b.Next,
		CreatedBy:    b.CreatedBy,
	}

	return &e
}

func BlockToModel(block entity.Block) *Block {
	model := Block{
		Hash:         block.Hash,
		PreviousHash: block.PreviousHash,
		Timestamp:    block.TimeStamp,
		Next:         block.Next,
		CreatedBy:    block.CreatedBy,
	}

	if block.Number != 0 {
		model.ID = uint(block.Number)
	}

	return &model
}
