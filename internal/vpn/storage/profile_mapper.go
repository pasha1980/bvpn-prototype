package storage

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func ProfileToStorage(profile entity.Profile) *ClientStorageFormat {
	return &ClientStorageFormat{
		Id:        profile.Id.String(),
		PrvKey:    fmt.Sprintf("%x", profile.PrvKey),
		PubKey:    fmt.Sprintf("%x", profile.PubKey),
		Client:    profile.Client,
		Price:     profile.Offer.Price,
		URL:       profile.Offer.URL,
		Timestamp: profile.Offer.Timestamp.Unix(),
	}
}

func StorageFormatToProfile(storageFormat ClientStorageFormat) (*entity.Profile, error) {
	id, err := uuid.Parse(storageFormat.Id)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id: id,
		Offer: block_data.Offer{
			URL:       storageFormat.URL,
			Price:     storageFormat.Price,
			Timestamp: time.Unix(storageFormat.Timestamp, 0),
		},
		PrvKey: []byte(storageFormat.PrvKey),
		PubKey: []byte(storageFormat.PubKey),
		Client: storageFormat.Client,
	}, nil
}
