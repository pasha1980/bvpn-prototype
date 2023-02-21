package data_validators

import (
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/hasher"
	"errors"
	"strconv"
)

func hashValidation(block entity.Block) error {
	hash := hasher.EncryptBlock(block)
	if string(hash) != block.Hash {
		return errors.New("Invalid hash on block #" + strconv.FormatUint(block.Number, 10)) // todo: Custom error
	}

	return nil
}
