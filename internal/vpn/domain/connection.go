package domain

import (
	"bvpn-prototype/internal/protocol/entity"
)

type Connection struct {
	Profile entity.Profile
	Traffic float64
}
