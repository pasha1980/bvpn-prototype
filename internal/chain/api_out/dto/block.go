package dto

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
)

type BlockDto struct {
	Number       uint64                   `json:"number"`
	Hash         string                   `json:"hash"`
	PreviousHash string                   `json:"previousHash"`
	Data         []block_data.ChainStored `json:"data"`
	TimeStamp    int64                    `json:"timeStamp"`
	Next         string                   `json:"next"`
	CreatedBy    string                   `json:"createdByy"`
}

func BlockToDto(block entity.Block) BlockDto {
	return BlockDto{
		Number:       block.Number,
		Hash:         block.Hash,
		PreviousHash: block.PreviousHash,
		Data:         block.Data,
		TimeStamp:    block.TimeStamp.Unix(),
		Next:         block.Next,
		CreatedBy:    block.CreatedBy,
	}
}
