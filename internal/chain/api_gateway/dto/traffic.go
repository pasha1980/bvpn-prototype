package dto

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type Traffic struct {
	ID        uuid.UUID `json:"id"`
	Sign      string    `json:"sign"`
	PubKey    string    `json:"pub"`
	Node      string    `json:"node"`
	Client    string    `json:"client"`
	Gb        float64   `json:"gb"`
	Timestamp int64     `json:"timestamp"`
}

func TrafficToDto(entity block_data.ChainStored) Traffic {
	traffic := entity.Data.(block_data.Traffic)
	return Traffic{
		ID:        entity.ID,
		Sign:      entity.Sign,
		PubKey:    entity.PubKey,
		Node:      traffic.Node,
		Client:    traffic.Client,
		Gb:        traffic.Bytes,
		Timestamp: traffic.Timestamp.Unix(),
	}
}
