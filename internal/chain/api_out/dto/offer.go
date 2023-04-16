package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Offer struct {
	ID        uuid.UUID `json:"id"`
	Sign      string    `json:"sign"`
	PubKey    string    `json:"pub"`
	URL       string    `json:"url"`
	Price     float64   `json:"price"`
	Timestamp int64     `json:"timestamp"`
}

func OfferToDto(entity block_data.ChainStored) Offer {
	offer := entity.Data.(block_data.Offer)
	return Offer{
		ID:        entity.ID,
		Sign:      entity.Sign,
		PubKey:    entity.PubKey,
		URL:       offer.URL,
		Price:     offer.Price,
		Timestamp: offer.Timestamp.Unix(),
	}
}
