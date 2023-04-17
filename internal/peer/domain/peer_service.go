package domain

import (
	"bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/peer/api_out"
	peer_errors "bvpn-prototype/internal/peer/errors"
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

	err = p.peerRepo.Save(peer)
	if err != nil {
		return errors.StorageError()
	}
	go api_out.AddMeTo(peer)
	return nil
}

func (p *PeerServiceImpl) GetPeers(except *entity.Node) []entity.Node {
	var result []entity.Node
	peers, err := p.peerRepo.GetAll()
	if err != nil {
		return nil
	}
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
			p.peerRepo.Remove(node) // todo
		}
	}
}

func (*PeerServiceImpl) validatePeer(peer entity.Node) error {
	ok := api_out.HealthCheck(peer)
	if !ok {
		return peer_errors.PeerNotAvailable(peer)
	}

	return protocol.ValidatePeer(peer)
}

func NewPeerService() (*PeerServiceImpl, error) {
	peerRepo, err := storage.NewPeerRepo()
	if err != nil {
		return nil, err
	}

	return &PeerServiceImpl{
		peerRepo: peerRepo,
	}, nil
}
