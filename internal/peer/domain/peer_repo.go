package domain

import "bvpn-prototype/internal/protocol/entity"

type PeerRepo interface {
	GetAll() []entity.Node
	Save(peer entity.Node)
	Remove(peer entity.Node)
}
