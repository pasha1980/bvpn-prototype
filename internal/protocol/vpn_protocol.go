package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/vpn_crypto"
	"github.com/google/uuid"
)

func GenerateVpnProfile(offer block_data.Offer, clientAddr string) entity.Profile {
	profile := entity.Profile{
		Id:     uuid.New(),
		Offer:  offer,
		Client: clientAddr,
	}

	profile.PrvKey, profile.PubKey = vpn_crypto.GeneratePair()
	return profile
}

func EncryptMessage(message []byte, profile entity.Profile) ([]byte, error) {
	return vpn_crypto.Encode(message, profile.PubKey)
}

func DecryptMessage(message []byte, profile entity.Profile) ([]byte, error) {
	return vpn_crypto.Decode(message, profile.PrvKey)
}
