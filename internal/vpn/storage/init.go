package storage

import (
	"bvpn-prototype/internal/storage/config"
	"bvpn-prototype/internal/storage/vpn_profile/errors"
	"os"
)

func InitStorage() error {
	dir := config.Get().StorageDirectory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return errors.FilesystemError
		}
	}

	return nil
}
