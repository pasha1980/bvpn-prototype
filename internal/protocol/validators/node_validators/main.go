package node_validators

import "bvpn-prototype/internal/protocol/entity"

var nodeValidationRules = []func(node entity.Node) error{
	ipCheck,
	pingCheck,
}

func GetValidationRules() []func(node entity.Node) error {
	return nodeValidationRules
}
