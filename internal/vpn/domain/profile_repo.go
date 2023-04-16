package domain

import "bvpn-prototype/internal/protocol/entity"

type ProfileRepo interface {
	Save(profile entity.Profile) (entity.Profile, error)
}
