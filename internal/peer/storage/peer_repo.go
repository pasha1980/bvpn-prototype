package storage

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/base64"
	"encoding/json"
	"os"
)

type PeerRepo struct {
	storage peerStorage
}

func (r *PeerRepo) GetAll() ([]entity.Node, error) {
	err := r.updateData()
	if err != nil {
		return nil, err
	}

	return r.storage.Data, nil
}

func (r *PeerRepo) Save(peer entity.Node) error {
	err := r.updateData()
	if err != nil {
		return err
	}
	if ok, _ := r.isExist(peer); !ok {
		r.storage.Data = append(r.storage.Data, peer)
		err = r.storeData()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PeerRepo) Remove(peer entity.Node) error {
	err := r.updateData()
	if err != nil {
		return err
	}
	if ok, _ := r.isExist(peer); !ok {
		data := r.storage.Data
		r.storage.Data = []entity.Node{}
		for _, existingPeer := range data {
			if peer != existingPeer {
				r.storage.Data = append(r.storage.Data, existingPeer)
			}
		}
		err = r.storeData()
		return err
	}

	return nil
}

func (r *PeerRepo) isExist(peer entity.Node) (bool, error) {
	err := r.updateData()
	if err != nil {
		return false, err
	}

	for _, stored := range r.storage.Data {
		if stored == peer {
			return true, nil
		}
	}

	return false, nil
}

func (r *PeerRepo) updateData() error {
	base64Encoded, err := os.ReadFile(config.Get().StorageDirectory + "/peers.bvpn")
	if err != nil {
		return r.storeData()
	}

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	if err != nil {
		return r.storeData()
	}

	err = json.Unmarshal(jsonEncoded, &r.storage)
	if err != nil {
		return r.storeData()
	}

	return nil
}

func (r *PeerRepo) storeData() error {
	jsonEncoded, _ := json.Marshal(r.storage)
	base64Encoded := base64.StdEncoding.EncodeToString(jsonEncoded)
	err := os.WriteFile(config.Get().StorageDirectory+"/peers.bvpn", []byte(base64Encoded), 0666)
	if err != nil {
		return err
	}

	return nil
}

func NewPeerRepo() (*PeerRepo, error) {
	repo := &PeerRepo{}
	err := repo.updateData()
	if err != nil {
		return nil, errors.StorageError(err.Error())
	}

	return repo, nil
}
