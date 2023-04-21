package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/protocol_error"
)

func typeValidation(block entity.Block, previousBlock *entity.Block) error {
	for _, data := range block.Data {
		switch data.Type {
		case block_data.TypeTransaction:
			tx := data.Data.(block_data.Transaction)
			if tx.From == "" || tx.Amount == 0 || tx.To == "" {
				return protocol_error.BlockValidationError("Invalid tx #"+data.ID.String(), block.Number)
			}
			break
		case block_data.TypeOffer:
			offer := data.Data.(block_data.Offer)
			if offer.Timestamp.Unix() == 0 || offer.URL == "" || offer.Price == 0 {
				return protocol_error.BlockValidationError("Invalid offer #"+data.ID.String(), block.Number)
			}
		case block_data.TypeConnectionBreak:
			cb := data.Data.(block_data.ConnectionBreak)
			if cb.Node == "" || cb.Timestamp.Unix() == 0 {
				return protocol_error.BlockValidationError("Invalid connection break #"+data.ID.String(), block.Number)
			}
		case block_data.TypeTraffic:
			traffic := data.Data.(block_data.Traffic)
			if traffic.Bytes == 0 || traffic.Node == "" || traffic.Timestamp.Unix() == 0 || traffic.Client == "" {
				return protocol_error.BlockValidationError("Invalid traffic #"+data.ID.String(), block.Number)
			}
		default:
			return protocol_error.BlockValidationError("Unsupported data type", block.Number)
		}
	}

	return nil
}
