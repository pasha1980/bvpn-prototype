package chain_stored_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	From       string    `json:"from" xml:"from" form:"from" query:"from"`
	To         string    `json:"to" xml:"to" form:"to" query:"to"`
	Amount     float64   `json:"amount" xml:"amount" form:"amount" query:"amount"`
	Commission float64   `json:"commission" xml:"commission" form:"commission" query:"commission"`
}

func (t *Transaction) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:   t.ID,
		Type: block_data.TypeTransaction,
		Data: block_data.Transaction{
			From:       t.From,
			To:         t.To,
			Amount:     t.Amount,
			Commission: t.Commission,
		},
	}
}
