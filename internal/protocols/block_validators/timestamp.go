package block_validators

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/protocol_error"
	"time"
)

func timestampValidation(block entity.Block, previousBlock *entity.Block) error {
	if block.TimeStamp.After(time.Now()) {
		return protocol_error.BlockValidationError("Invalid timestamp", block.Number)
	}

	if previousBlock.TimeStamp.Add(10 * time.Second).After(block.TimeStamp) {
		return protocol_error.BlockValidationError("Invalid timestamp", block.Number)
	}

	return nil
}
