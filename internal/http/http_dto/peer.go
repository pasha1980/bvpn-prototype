package http_dto

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
)

type PeerDto struct {
	URL string `json:"url"`
	// todo
}

func (d *PeerDto) ToEntity() entity.Node {
	return entity.Node{
		URL: d.URL,
	}
}

func MakeMeDTO() PeerDto {
	return PeerDto{
		URL: config.Get().URL,
	}
}
