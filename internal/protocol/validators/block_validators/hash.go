package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/hasher"
	"bvpn-prototype/internal/protocol/protocol_error"
)

func hashValidation(block entity.Block, previousBlock *entity.Block) error {
	hash := hasher.EncryptBlock(block)
	if hash != block.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	if block.PreviousHash != previousBlock.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	return nil
}
