package entity

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"time"
)

type Block struct {
	Number       uint64
	Hash         string
	PreviousHash string
	Data         []block_data.ChainStored
	TimeStamp    time.Time
	CreatedBy    string
	Next         string
}

func (b *Block) IsInitial() bool {
	return b.Number == 1
}
