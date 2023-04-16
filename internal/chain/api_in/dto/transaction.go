package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Transaction struct {
	ID     uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Sign   string    `json:"sign" xml:"sign" form:"sign" query:"sign"`
	PubKey string    `json:"pub" xml:"pub" form:"pub" query:"pub"`
	From   string    `json:"from" xml:"from" form:"from" query:"from"`
	To     string    `json:"to" xml:"to" form:"to" query:"to"`
	Amount float64   `json:"amount" xml:"amount" form:"amount" query:"amount"`
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
