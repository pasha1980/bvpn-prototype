package models

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type Offer struct {
	ID      uint `gorm:"PRIMARY_KEY"`
	Ref     string
	BlockID uint `gorm:"index"`
	Sign    string
	PubKey  string
	//Node    string
	Price     float64
	Timestamp int64
}

func (o Offer) TableName() string {
	return "offer"
}

func (o *Offer) ToChainStored() block_data.ChainStored {
	id, _ := uuid.FromBytes([]byte(o.Ref))
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeOffer,
		Sign:   o.Sign,
		PubKey: o.PubKey,
		Data: block_data.Offer{
			//Node:  o.Node,
			Price:     o.Price,
			Timestamp: o.Timestamp,
		},
	}
}

func OfferToModel(data block_data.ChainStored, blockId uint) Offer {
	offer := data.Data.(block_data.Offer)
	return Offer{
		Ref:     data.ID.String(),
		BlockID: blockId,
		Sign:    data.Sign,
		PubKey:  data.PubKey,
		//Node:    offer.Node,
		Price:     offer.Price,
		Timestamp: offer.Timestamp,
	}
}
