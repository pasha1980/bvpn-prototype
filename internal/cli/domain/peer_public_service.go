package domain

import "bvpn-prototype/internal/protocol/entity"

type PeerPublicService interface {
	AddPeer(peer entity.Node) error
}
