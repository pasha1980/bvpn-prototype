package domain

import (
	"bvpn-prototype/internal/protocol/entity"
	"github.com/google/uuid"
)

type ProfileRepo interface {
	Save(profile entity.Profile) (*entity.Profile, error)
	Get(id uuid.UUID) (*entity.Profile, error)
	Remove(id uuid.UUID) error
	IsExist(id uuid.UUID) (bool, error)
}
