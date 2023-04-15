package block_validators

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/protocol_error"
	"bvpn-prototype/internal/protocols/signer"
)

func signValidation(block entity.Block, previousBlock *entity.Block) error {
	for _, data := range block.Data {
		if !signer.Validate(&data) {
			return protocol_error.BlockValidationError("Invalid signature on data #"+data.ID.String(), block.Number)
		}
	}

	return nil
}
