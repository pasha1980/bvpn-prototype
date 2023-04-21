package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/hasher"
	"bvpn-prototype/internal/protocol/interfaces"
	"bvpn-prototype/internal/protocol/protocol_error"
)

func hashValidation(block entity.Block, reader interfaces.ChainReader) error {
	hash := hasher.EncryptBlock(block)
	if hash != block.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	previousBlock := reader.Previous(block.Number)
	if block.PreviousHash != previousBlock.Hash {
		return protocol_error.BlockValidationError("Invalid hash", block.Number)
	}

	return nil
}
