package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
	"time"
)

type ConnectionBreak struct {
	ID        uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Sign      string    `json:"sign" xml:"sign" form:"sign" query:"sign"`
	PubKey    string    `json:"pub" xml:"pub" form:"pub" query:"pub"`
	Node      string    `json:"node" xml:"node" form:"node" query:"node"`
	Timestamp int64     `json:"active" xml:"active" form:"active" query:"active"`
}

func (b *ConnectionBreak) ToEntity() block_data.ChainStored {
	timestamp := time.Unix(b.Timestamp, 0)
	return block_data.ChainStored{
		ID:     b.ID,
		Sign:   b.Sign,
		PubKey: b.PubKey,
		Type:   block_data.TypeTransaction,
		Data: block_data.ConnectionBreak{
			Node:      b.Node,
			Timestamp: timestamp,
		},
	}
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
