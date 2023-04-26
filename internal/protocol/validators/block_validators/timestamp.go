package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/params"
	"time"
)

func timestampValidation(block entity.Block, reader interfaces.ChainReader) error {
	if block.TimeStamp.After(time.Now()) {
		return errors.BlockValidationError("Invalid timestamp", block.Number)
	}

	previousBlock := reader.Previous(block.Number)
	if previousBlock.TimeStamp.Add(params.TimeToWaitNextBlock).After(block.TimeStamp) {
		return errors.BlockValidationError("Invalid timestamp", block.Number)
	}

	return nil
}
