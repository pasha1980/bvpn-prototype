package entity

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"time"
)

const InitialBlockTimestamp = "2023-02-25 00:00:00"
const InitialBlockPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

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
