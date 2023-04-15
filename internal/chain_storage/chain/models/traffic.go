package models

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
	"time"
)

type Traffic struct {
	ID        uint `gorm:"PRIMARY_KEY"`
	Ref       string
	BlockID   uint `gorm:"index"`
	Sign      string
	PubKey    string
	Node      string
	Client    string
	Gb        float64
	Timestamp int64
}

func (t Traffic) TableName() string {
	return "traffic"
}

func (t *Traffic) ToChainStored() block_data.ChainStored {
	id, _ := uuid.FromBytes([]byte(t.Ref))
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeTraffic,
		Sign:   t.Sign,
		PubKey: t.PubKey,
		Data: block_data.Traffic{
			Node:      t.Node,
			Client:    t.Client,
			Gb:        t.Gb,
			Timestamp: time.Unix(t.Timestamp, 0),
		},
	}
}

func TrafficToModel(data block_data.ChainStored, blockId uint) Traffic {
	traffic := data.Data.(block_data.Traffic)
	return Traffic{
		Ref:       data.ID.String(),
		BlockID:   blockId,
		Sign:      data.Sign,
		PubKey:    data.PubKey,
		Node:      traffic.Node,
		Client:    traffic.Client,
		Gb:        traffic.Gb,
		Timestamp: traffic.Timestamp.Unix(),
	}
}
