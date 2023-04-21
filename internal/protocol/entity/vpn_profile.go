package entity

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Profile struct {
	Id     uuid.UUID
	Client string
	Offer  block_data.Offer
	PrvKey []byte
	PubKey []byte
}
