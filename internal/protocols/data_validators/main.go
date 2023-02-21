package data_validators

import "bvpn-prototype/internal/protocols/entity"

var validationRules = []func(block entity.Block) error{
	hashValidation,
	timestampValidation,
}

func GetValidationRules() []func(block entity.Block) error {
	return validationRules
}
