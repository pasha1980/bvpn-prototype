package protocols

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/repo"
	"bvpn-prototype/internal/protocols/validators/node_validators"
	"bvpn-prototype/internal/storage/peer"
)

type PeerProtocol struct {
	repo repo.PeerStorageRepo
}

func (p *PeerProtocol) GetPeers() []entity.Node {
	return p.repo.GetAll()
}

func (p *PeerProtocol) AddNewPeer(peer entity.Node) error {
	err := p.ValidatePeer(peer)
	if err != nil {
		return err
	}

	p.repo.Save(peer)
	return nil
}

func (p *PeerProtocol) ValidatePeer(peer entity.Node) error {
	var err error

	for _, validator := range node_validators.GetValidationRules() {
		err = validator(peer)
	}

	return err
}

func (p *PeerProtocol) CheckPeers() {
	peers := p.GetPeers()
	for _, node := range peers {
		err := p.ValidatePeer(node)
		if err != nil {
			p.repo.Remove(node)
		}
	}
}

func GetPeerProtocol() *PeerProtocol {
	return &PeerProtocol{
		repo: peer.NewPeerRepo(),
	}
}
