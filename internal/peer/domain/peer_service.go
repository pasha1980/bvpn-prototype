package domain

import (
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/errors"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity"
)

type PeerService interface {
	CheckPeers()
	AddPeer(peer entity.Node) error
}

type PeerServiceImpl struct {
}

func (p *PeerServiceImpl) AddPeer(peer entity.Node) error {
	err := p.validatePeer(peer)
	if err != nil {
		return err
	}

	err = p.repo().Save(peer)
	if err != nil {
		return errors.StorageError()
	}
	go p.gateway().Peer(peer).AddMe()
	return nil
}

func (p *PeerServiceImpl) GetPeers(except *entity.Node) []entity.Node {
	var result []entity.Node
	peers, err := p.repo().GetAll()
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
			p.repo().Remove(node)
		}
	}
}

func (p *PeerServiceImpl) validatePeer(peer entity.Node) error {
	return protocol.ValidatePeer(peer)
}

func (*PeerServiceImpl) repo() PeerRepo {
	return di.Get("peer_repo").(PeerRepo)
}

func (*PeerServiceImpl) gateway() PeerApiGateway {
	return di.Get("peer_api_gateway").(PeerApiGateway)
}

func NewPeerService() (*PeerServiceImpl, error) {
	return &PeerServiceImpl{}, nil
}
