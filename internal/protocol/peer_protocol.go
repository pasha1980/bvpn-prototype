package protocol

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/validators/node_validators"
)

func ValidatePeer(peer entity.Node) error {
	var err error

	for _, validator := range node_validators.GetValidationRules() {
		err = validator(peer)
	}

	return err
}
