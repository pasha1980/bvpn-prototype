package storage

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/pem"
	"os"
)

type ProfileRepo struct {
	dir string
}

func (r *ProfileRepo) Save(profile entity.Profile) (entity.Profile, error) {
	keyFile := r.dir + "/" + profile.Id + "/prv.pem"
	var file *os.File
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		file, err = os.Create(keyFile)
		if err != nil {
			// todo
		}
	} else {
		file, err = os.Open(keyFile)
		if err != nil {
			// todo
		}
	}

	pem.Encode(file, &pem.Block{
		Type:    "BVPN KEY",
		Headers: nil,
		Bytes:   profile.PrvKey,
	})

	return profile, nil
}

func (r *ProfileRepo) Get(id string) (*entity.Profile, error) {
	// tdodo
	return nil, nil
}

func NewProfileRepo() (*ProfileRepo, error) {
	baseDir := config.Get().StorageDirectory
	dir := baseDir + "/profiles"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			// todo
		}
	}
	return &ProfileRepo{
		dir: dir,
	}, nil
}
