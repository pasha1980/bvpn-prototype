package block_validators

import (
	"bvpn-prototype/internal/protocols/entity"
)

var blockValidationRules = []func(block entity.Block, previousBlock *entity.Block) error{
	hashValidation,
	timestampValidation,
	signValidation,
	offerValidation,
}

func GetValidationRules() []func(block entity.Block, previousBlock *entity.Block) error {
	return blockValidationRules
}
