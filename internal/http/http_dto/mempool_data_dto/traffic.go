package mempool_data_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
	"time"
)

type Traffic struct {
	ID        uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Node      string    `json:"node" xml:"node" form:"node" query:"node"`
	Client    string    `json:"client" xml:"client" form:"client" query:"client"`
	Gb        float64   `json:"gb" xml:"gb" form:"gb" query:"gb"`
	Timestamp time.Time `json:"timestamp" xml:"timestamp" form:"timestamp" query:"timestamp"`
}

func (t *Traffic) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:   t.ID,
		Type: block_data.TypeTransaction,
		Data: block_data.Traffic{
			Node:      t.Node,
			Client:    t.Client,
			Gb:        t.Gb,
			Timestamp: t.Timestamp,
		},
	}
}

func TrafficToDto(entity block_data.ChainStored) Traffic {
	traffic := entity.Data.(block_data.Traffic)
	return Traffic{
		ID:        entity.ID,
		Node:      traffic.Node,
		Client:    traffic.Client,
		Gb:        traffic.Gb,
		Timestamp: traffic.Timestamp,
	}
}
