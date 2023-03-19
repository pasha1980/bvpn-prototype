package block_data

import "github.com/google/uuid"

type ChainStored struct {
	ID     uuid.UUID
	Type   StoredEntityType
	Sign   string
	PubKey string
	Data   any
}

type StoredEntityType uint8
