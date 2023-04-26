package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/interfaces"
)

var blockValidationRules = []func(block entity.Block, reader interfaces.ChainReader) error{
	hashValidation,
	timestampValidation,
	signValidation,
	typeValidation,
	offerValidation,
	dataSizeValidation,
}

func GetValidationRules() []func(block entity.Block, reader interfaces.ChainReader) error {
	return blockValidationRules
}
