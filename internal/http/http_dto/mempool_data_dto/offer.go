package mempool_data_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type Offer struct {
	ID     uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Sign   string    `json:"sign" xml:"sign" form:"sign" query:"sign"`
	PubKey string    `json:"pub" xml:"pub" form:"pub" query:"pub"`
	//Node   string    `json:"node" xml:"node" form:"node" query:"node"`
	Price     float64 `json:"price" xml:"price" form:"price" query:"price"`
	Timestamp int64   `json:"timestamp" xml:"timestamp" form:"timestamp" query:"timestamp"`
}

func (o *Offer) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:     o.ID,
		Sign:   o.Sign,
		PubKey: o.PubKey,
		Type:   block_data.TypeTransaction,
		Data: block_data.Offer{
			//Node:  o.Node,
			Price:     o.Price,
			Timestamp: o.Timestamp,
		},
	}
}

func OfferToDto(entity block_data.ChainStored) Offer {
	offer := entity.Data.(block_data.Offer)
	return Offer{
		ID:     entity.ID,
		Sign:   entity.Sign,
		PubKey: entity.PubKey,
		//Node:   offer.Node,
		Price:     offer.Price,
		Timestamp: offer.Timestamp,
	}
}
