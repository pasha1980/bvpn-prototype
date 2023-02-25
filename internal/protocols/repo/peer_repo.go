package repo

import "bvpn-prototype/internal/protocols/entity"

type PeerStorageRepo interface {
	GetAll() []entity.Node
	Save(peer entity.Node)
	Remove(peer entity.Node)
}
