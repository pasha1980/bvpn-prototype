package block_data

import "github.com/google/uuid"

type ChainStored struct {
	ID     uuid.UUID        `json:"id"`
	Type   StoredEntityType `json:"type"`
	Sign   string           `json:"sign"`
	PubKey string           `json:"pubKey"`
	Data   any              `json:"data"`
}

type StoredEntityType uint8
