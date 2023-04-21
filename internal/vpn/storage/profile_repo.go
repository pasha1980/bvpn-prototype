package storage

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity"
	"github.com/google/uuid"
	"os"
)

const keyFileName = "client.bvpn"

type ProfileRepo struct {
	dir string
}

func (r *ProfileRepo) Save(profile entity.Profile) (*entity.Profile, error) {
	keyFile := r.dir + "/" + profile.Id.String() + "/" + keyFileName
	var file *os.File
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		file, err = os.Create(keyFile)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(keyFile)
		if err != nil {
			return nil, err
		}
	}

	err := ProfileToStorage(profile).Write(file)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepo) Get(id uuid.UUID) (*entity.Profile, error) {
	if ok, err := r.IsExist(id); !ok {
		return nil, err
	}

	keyFile := r.dir + "/" + id.String() + "/" + keyFileName
	file, err := os.Open(keyFile)
	if err != nil {
		return nil, err
	}

	format, err := Read(file)
	if err != nil {
		return nil, err
	}

	profile, err := StorageFormatToProfile(*format)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepo) Remove(id uuid.UUID) error {
	// todo
	return nil
}

func (r *ProfileRepo) IsExist(id uuid.UUID) (bool, error) {
	keyFile := r.dir + "/" + id.String() + "/" + keyFileName
	_, err := os.Stat(keyFile)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
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
