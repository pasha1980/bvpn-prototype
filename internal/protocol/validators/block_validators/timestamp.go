package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/protocol_error"
	"time"
)

func timestampValidation(block entity.Block, reader interfaces.ChainReader) error {
	if block.TimeStamp.After(time.Now()) {
		return protocol_error.BlockValidationError("Invalid timestamp", block.Number)
	}

	previousBlock := reader.Previous(block.Number)
	if previousBlock.TimeStamp.Add(10 * time.Second).After(block.TimeStamp) {
		return protocol_error.BlockValidationError("Invalid timestamp", block.Number)
	}

	return nil
}
