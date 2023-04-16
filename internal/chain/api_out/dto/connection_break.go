package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type ConnectionBreak struct {
	ID        uuid.UUID `json:"id"`
	Sign      string    `json:"sign"`
	PubKey    string    `json:"pub"`
	Node      string    `json:"node"`
	Timestamp int64     `json:"active"`
}

func ConnectionBreakToDto(entity block_data.ChainStored) ConnectionBreak {
	cb := entity.Data.(block_data.ConnectionBreak)
	return ConnectionBreak{
		ID:        entity.ID,
		Sign:      entity.Sign,
		PubKey:    entity.PubKey,
		Node:      cb.Node,
		Timestamp: cb.Timestamp.Unix(),
	}
}
