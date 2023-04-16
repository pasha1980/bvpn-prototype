package storage

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/base64"
	"encoding/json"
	"os"
)

type PeerRepo struct {
	storage peerStorage
}

func (r *PeerRepo) GetAll() []entity.Node {
	r.updateData()
	return r.storage.Data
}

func (r *PeerRepo) Save(peer entity.Node) {
	r.updateData()
	if !r.isExist(peer) {
		r.storage.Data = append(r.storage.Data, peer)
		r.storeData()
	}
}

func (r *PeerRepo) Remove(peer entity.Node) {
	r.updateData()
	if r.isExist(peer) {
		data := r.storage.Data
		r.storage.Data = []entity.Node{}
		for _, existingPeer := range data {
			if peer != existingPeer {
				r.storage.Data = append(r.storage.Data, existingPeer)
			}
		}
		r.storeData()
	}
}

func (r *PeerRepo) isExist(peer entity.Node) bool {
	r.updateData()
	for _, stored := range r.storage.Data {
		if stored == peer {
			return true
		}
	}

	return false
}

func (r *PeerRepo) updateData() {
	base64Encoded, err := os.ReadFile(config.Get().StorageDirectory + "/peers.bvpn")
	if err != nil {
		r.storeData()
	}

	jsonEncoded, err := base64.StdEncoding.DecodeString(string(base64Encoded))
	if err != nil {
		r.storeData()
	}

	err = json.Unmarshal(jsonEncoded, &r.storage)
	if err != nil {
		r.storeData()
	}
}

func (r *PeerRepo) storeData() {
	jsonEncoded, _ := json.Marshal(r.storage)
	base64Encoded := base64.StdEncoding.EncodeToString(jsonEncoded)
	err := os.WriteFile(config.Get().StorageDirectory+"/peers.bvpn", []byte(base64Encoded), 0666)
	if err != nil {
		// todo: what to do
	}
}

func NewPeerRepo() *PeerRepo {
	repo := &PeerRepo{}
	repo.updateData()
	return repo
}
