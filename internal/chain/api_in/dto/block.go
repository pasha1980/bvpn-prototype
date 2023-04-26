package dto

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"time"
)

type BlockDto struct {
	Number       uint64                   `json:"number"`
	Hash         string                   `json:"hash"`
	PreviousHash string                   `json:"previousHash"`
	Data         []block_data.ChainStored `json:"data" validate:"dive"`
	TimeStamp    int64                    `json:"timeStamp" `
	Next         string                   `json:"next"`
	CreatedBy    string                   `json:"createdBy"`
}

func (d *BlockDto) ToEntity() entity.Block {
	var data []block_data.ChainStored
	for _, dto := range d.Data {
		switch dto.Type {
		case block_data.TypeTransaction:
			data = append(data, block_data.ChainStored{
				ID:   dto.ID,
				Type: dto.Type,
				Data: dto.Data.(block_data.Transaction),
			})
			break
		case block_data.TypeOffer:
			data = append(data, block_data.ChainStored{
				ID:   dto.ID,
				Type: dto.Type,
				Data: dto.Data.(block_data.Offer),
			})
			break
		case block_data.TypeConnectionBreak:
			data = append(data, block_data.ChainStored{
				ID:   dto.ID,
				Type: dto.Type,
				Data: dto.Data.(block_data.ConnectionBreak),
			})
			break
		case block_data.TypeTraffic:
			data = append(data, block_data.ChainStored{
				ID:   dto.ID,
				Type: dto.Type,
				Data: dto.Data.(block_data.Traffic),
			})
			break
		default:
			data = append(data, block_data.ChainStored{
				ID:   dto.ID,
				Type: dto.Type,
				Data: dto.Data.(map[string]any),
			})
			break
		}
	}

	return entity.Block{
		Number:       d.Number,
		Hash:         d.Hash,
		PreviousHash: d.PreviousHash,
		Data:         data,
		TimeStamp:    time.Unix(d.TimeStamp, 0),
		Next:         d.Next,
		CreatedBy:    d.CreatedBy,
	}
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
