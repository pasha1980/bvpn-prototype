package storage

import (
	"bvpn-prototype/internal/storage/config"
	"bvpn-prototype/internal/storage/vpn_profile/errors"
	"encoding/pem"
	"os"
)

type ProfileRepo struct {
	dir string
}

func (r *ProfileRepo) Save(profile Profile) (Profile, error) {
	keyFile := r.dir + "/" + profile.id + "/prv.pem"
	var file *os.File
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		file, err = os.Create(keyFile)
		if err != nil {
			return profile, errors.FilesystemError
		}
	} else {
		file, err = os.Open(keyFile)
		if err != nil {
			return profile, errors.FilesystemError
		}
	}

	pem.Encode(file, &pem.Block{
		Type:    "BVPN KEY",
		Headers: nil,
		Bytes:   profile.PrvKey,
	})

	return profile, nil
}

func (r *ProfileRepo) Get(id string) (*Profile, error) {
	// tdodo
	return nil, nil
}

func NewRepo() *ProfileRepo {
	baseDir := config.Get().StorageDirectory
	return &ProfileRepo{
		dir: baseDir + "/profiles",
	}
}
