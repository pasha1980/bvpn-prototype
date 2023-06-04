package domain

import "bvpn-prototype/internal/protocol/entity"

type PeerApiGateway interface {
	Peer(node entity.Node) PeerApiGateway
	AddMe()
	HealthCheck() bool
}
