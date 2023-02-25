package block_validators

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/hasher"
	"bvpn-prototype/internal/protocols/protocol_error"
)

func hashValidation(block entity.Block, previousBlock *entity.Block) error {
	hash := hasher.EncryptBlock(block)
	if string(hash) != block.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	if block.PreviousHash != previousBlock.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	return nil
}
