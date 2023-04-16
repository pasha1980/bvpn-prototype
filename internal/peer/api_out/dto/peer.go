package dto

import (
	"bvpn-prototype/internal/infrastructure/config"
)

type PeerDto struct {
	URL string `json:"url"`
	// todo
}

func MakeMeDTO() PeerDto {
	return PeerDto{
		URL: config.Get().URL,
	}
}
