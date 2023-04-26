package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/params"
	"bvpn-prototype/utils"
)

func dataSizeValidation(block entity.Block, reader interfaces.ChainReader) error {
	if utils.SizeOf(block.Data) > params.BlockCapacity {
		return errors.BlockValidationError("Data is out of available capacity", block.Number)
	}

	return nil
}
