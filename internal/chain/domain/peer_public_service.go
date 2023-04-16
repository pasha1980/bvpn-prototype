package domain

import "bvpn-prototype/internal/protocol/entity"

type PeerPublicService interface {
	GetPeers(except *entity.Node) []entity.Node
}
