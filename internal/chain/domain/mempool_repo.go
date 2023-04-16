package domain

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"github.com/google/uuid"
)

type MempoolRepository interface {
	AddNewElement(element block_data.ChainStored)
	GetElements(size int) []block_data.ChainStored
	IsExist(uuid uuid.UUID) bool
	RemoveByIndex(index uuid.UUID)
}
