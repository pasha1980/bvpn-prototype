package dto

import (
	"bvpn-prototype/internal/protocol/entity"
)

type PeerDto struct {
	URL string `json:"url" validate:"url"`
	// todo
}

func (d *PeerDto) ToEntity() entity.Node {
	return entity.Node{
		URL: d.URL,
	}
}
