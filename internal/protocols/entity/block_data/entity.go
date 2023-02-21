package block_data

import "github.com/google/uuid"

type ChainStored struct {
	ID   uuid.UUID        `json:"id"`
	Type StoredEntityType `json:"type"`
	Data any              `json:"data"`
}

type StoredEntityType uint8
