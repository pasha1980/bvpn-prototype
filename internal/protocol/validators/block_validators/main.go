package block_validators

import (
	"bvpn-prototype/internal/protocol/entity"
)

var blockValidationRules = []func(block entity.Block, previousBlock *entity.Block) error{
	hashValidation,
	timestampValidation,
	signValidation,
	typeValidation,
	offerValidation,
}

func GetValidationRules() []func(block entity.Block, previousBlock *entity.Block) error {
	return blockValidationRules
}
