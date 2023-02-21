package chain_stored_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type Offer struct {
	ID    uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Node  string    `json:"node" xml:"node" form:"node" query:"node"`
	Price float64   `json:"price" xml:"price" form:"price" query:"price"`
}

func (o *Offer) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:   o.ID,
		Type: block_data.TypeTransaction,
		Data: block_data.Offer{
			Node:  o.Node,
			Price: o.Price,
		},
	}
}
