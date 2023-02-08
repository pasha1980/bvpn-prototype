package chain_domain

import (
	"bvpn-prototype/internal/chain/chain_domain/entity"
	"time"
)

type Block struct {
	Number       uint64
	Hash         string
	PreviousHash string
	Data         []entity.ChainStored
	TimeStamp    time.Time
}
