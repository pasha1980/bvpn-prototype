package peer

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
	"encoding/base64"
	"encoding/json"
	"os"
)

type PeerRepo struct {
	storage PeerStorage
}

func (r *PeerRepo) GetAll() []entity.Node {
	r.updateData()
	return r.storage.data
}

func (r *PeerRepo) Save(peer entity.Node) {
	r.updateData()
	if !r.isExist(peer) {
		r.storage.data = append(r.storage.data, peer)
		r.storeData()
	}
}

func (r *PeerRepo) Remove(peer entity.Node) {
	r.updateData()
	if r.isExist(peer) {
		data := r.storage.data
		r.storage.data = []entity.Node{}
		for _, existingPeer := range data {
			if peer != existingPeer {
				r.storage.data = append(r.storage.data, existingPeer)
			}
		}
		r.storeData()
	}
}

func (r *PeerRepo) isExist(peer entity.Node) bool {
	r.updateData()
	for _, stored := range r.storage.data {
		if stored == peer {
			return true
		}
	}

	return false
}

func (r *PeerRepo) updateData() {
	base64Encoded, err := os.ReadFile(config.Get().StorageDirectory + "/peers.txt")
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
	err := os.WriteFile(config.Get().StorageDirectory+"/peers.txt", []byte(base64Encoded), 0666)
	if err != nil {
		// todo: what to do
	}
}

func NewPeerRepo() *PeerRepo {
	repo := &PeerRepo{}
	repo.updateData()
	repo.storeData()
	return repo
}
