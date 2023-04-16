package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
	"time"
)

type Traffic struct {
	ID        uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Sign      string    `json:"sign" xml:"sign" form:"sign" query:"sign"`
	PubKey    string    `json:"pub" xml:"pub" form:"pub" query:"pub"`
	Node      string    `json:"node" xml:"node" form:"node" query:"node"`
	Client    string    `json:"client" xml:"client" form:"client" query:"client"`
	Gb        float64   `json:"gb" xml:"gb" form:"gb" query:"gb"`
	Timestamp int64     `json:"timestamp" xml:"timestamp" form:"timestamp" query:"timestamp"`
}

func (t *Traffic) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:     t.ID,
		Sign:   t.Sign,
		PubKey: t.PubKey,
		Type:   block_data.TypeTransaction,
		Data: block_data.Traffic{
			Node:      t.Node,
			Client:    t.Client,
			Gb:        t.Gb,
			Timestamp: time.Unix(t.Timestamp, 0),
		},
	}
}

func TrafficToDto(entity block_data.ChainStored) Traffic {
	traffic := entity.Data.(block_data.Traffic)
	return Traffic{
		ID:        entity.ID,
		Sign:      entity.Sign,
		PubKey:    entity.PubKey,
		Node:      traffic.Node,
		Client:    traffic.Client,
		Gb:        traffic.Gb,
		Timestamp: traffic.Timestamp.Unix(),
	}
}
