package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
	"bvpn-prototype/internal/protocol/interfaces"
)

func creatorValidation(block entity.Block, reader interfaces.ChainReader) error {
	previous := reader.Previous(block.Number)
	if previous.Next != block.CreatedBy {
		return errors.BlockValidationError("Wrong creation queue", block.Number)
	}

	if block.Next == "" {
		return errors.BlockValidationError("Next is not chosen", block.Number)
	}

	// todo: what if chain stopped
	// todo: what if next does not exist

	return nil
}
