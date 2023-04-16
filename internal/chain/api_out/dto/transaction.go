package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Transaction struct {
	ID     uuid.UUID `json:"id"`
	Sign   string    `json:"sign"`
	PubKey string    `json:"pub"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
}

func (t *Transaction) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:   t.ID,
		Type: block_data.TypeTransaction,
		Data: block_data.Transaction{
			From:   t.From,
			To:     t.To,
			Amount: t.Amount,
		},
	}
}

func TransactionToDto(entity block_data.ChainStored) Transaction {
	tx := entity.Data.(block_data.Transaction)
	return Transaction{
		ID:     entity.ID,
		Sign:   entity.Sign,
		PubKey: entity.PubKey,
		From:   tx.From,
		To:     tx.To,
		Amount: tx.Amount,
	}
}
