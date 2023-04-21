package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/vpn_crypto"
	"github.com/google/uuid"
)

func GenerateVpnProfile(offer block_data.Offer) entity.Profile {
	profile := entity.Profile{
		Id:    uuid.New(),
		Offer: offer,
	}

	profile.PrvKey, profile.PubKey = vpn_crypto.GeneratePair()
	return profile
}
