package domain

import "bvpn-prototype/internal/protocol/entity"

type PeerRepo interface {
	GetAll() ([]entity.Node, error)
	Save(peer entity.Node) error
	Remove(peer entity.Node) error
}
