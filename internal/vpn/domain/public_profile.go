package domain

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity"
	"github.com/google/uuid"
)

type PublicProfile struct {
	Id     uuid.UUID
	Proto  string
	Port   string
	PubKey string
}

func ProfileToPub(profile entity.Profile) PublicProfile {
	return PublicProfile{
		Id:     profile.Id,
		Proto:  config.Get().VpnProto,
		Port:   config.Get().VpnPort,
		PubKey: string(profile.PubKey),
	}
}
