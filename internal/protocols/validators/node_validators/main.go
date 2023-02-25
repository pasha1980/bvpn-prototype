package node_validators

import "bvpn-prototype/internal/protocols/entity"

var nodeValidationRules = []func(node entity.Node) error{}

func GetValidationRules() []func(node entity.Node) error {
	return nodeValidationRules
}
