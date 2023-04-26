package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/signer"
)

func signValidation(block entity.Block, reader interfaces.ChainReader) error {
	for _, data := range block.Data {
		if !signer.Validate(&data) {
			return errors.BlockValidationError("Invalid signature on data #"+data.ID.String(), block.Number)
		}
	}

	return nil
}
