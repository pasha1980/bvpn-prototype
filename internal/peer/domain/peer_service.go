package domain

import (
	"bvpn-prototype/internal/peer/api_out"
	"bvpn-prototype/internal/peer/errors"
	"bvpn-prototype/internal/peer/storage"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity"
)

type PeerService interface {
	CheckPeers()
	AddPeer(peer entity.Node) error
}

type PeerServiceImpl struct {
	peerRepo PeerRepo
}

func (p *PeerServiceImpl) AddPeer(peer entity.Node) error {
	err := p.validatePeer(peer)
	if err != nil {
		return err
	}

	p.peerRepo.Save(peer)
	go api_out.AddMeTo(peer)
	return nil
}

func (p *PeerServiceImpl) GetPeers(except *entity.Node) []entity.Node {
	var result []entity.Node
	peers := p.peerRepo.GetAll()
	if except != nil {
		for _, node := range peers {
			if except.IP != node.IP {
				result = append(result, node)
			}
		}
	} else {
		result = peers
	}

	return result
}

func (p *PeerServiceImpl) CheckPeers() {
	peers := p.GetPeers(nil)
	for _, node := range peers {
		err := p.validatePeer(node)
		if err != nil {
			p.peerRepo.Remove(node)
		}
	}
}

func (*PeerServiceImpl) validatePeer(peer entity.Node) error {
	ok := api_out.HealthCheck(peer)
	if !ok {
		return errors.PeerNotAvailable(peer)
	}

	return protocol.ValidatePeer(peer)
}

func NewPeerService() (*PeerServiceImpl, error) {
	return &PeerServiceImpl{
		peerRepo: storage.NewPeerRepo(),
	}, nil
}
