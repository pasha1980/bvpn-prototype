package models

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
	"time"
)

type ConnectionBreak struct {
	ID        uint `gorm:"PRIMARY_KEY"`
	Ref       string
	BlockID   uint `gorm:"index"`
	Sign      string
	PubKey    string
	Node      string
	Timestamp int64
}

func (b ConnectionBreak) TableName() string {
	return "connection_break"
}

func (b *ConnectionBreak) ToChainStored() block_data.ChainStored {
	id, _ := uuid.FromBytes([]byte(b.Ref))
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeConnectionBreak,
		Sign:   b.Sign,
		PubKey: b.PubKey,
		Data: block_data.ConnectionBreak{
			Node:      b.Node,
			Timestamp: time.Unix(b.Timestamp, 0),
		},
	}
}

func NodeStatusToModel(data block_data.ChainStored, blockId uint) ConnectionBreak {
	status := data.Data.(block_data.ConnectionBreak)
	return ConnectionBreak{
		Ref:       data.ID.String(),
		BlockID:   blockId,
		Sign:      data.Sign,
		PubKey:    data.PubKey,
		Node:      status.Node,
		Timestamp: status.Timestamp.Unix(),
	}
}
