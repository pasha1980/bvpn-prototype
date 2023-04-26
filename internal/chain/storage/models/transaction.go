package models

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Transaction struct {
	ID      uint `gorm:"PRIMARY_KEY"`
	Ref     string
	BlockID uint `gorm:"index"`
	Sign    string
	PubKey  string
	From    string `gorm:"index"`
	To      string `gorm:"index"`
	Amount  float64
}

func (t Transaction) TableName() string {
	return "tx"
}

func (t *Transaction) ToChainStored() block_data.ChainStored {
	id, _ := uuid.Parse(t.Ref)
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeTransaction,
		Sign:   t.Sign,
		PubKey: t.PubKey,
		Data: block_data.Transaction{
			To:     t.To,
			From:   t.From,
			Amount: t.Amount,
		},
	}
}

func TxToModel(data block_data.ChainStored, blockId uint) Transaction {
	tx := data.Data.(block_data.Transaction)
	return Transaction{
		Ref:     data.ID.String(),
		BlockID: blockId,
		Sign:    data.Sign,
		PubKey:  data.PubKey,
		From:    tx.From,
		To:      tx.To,
		Amount:  tx.Amount,
	}
}
