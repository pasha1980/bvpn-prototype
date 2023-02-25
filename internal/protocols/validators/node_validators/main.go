package node_validators

import "bvpn-prototype/internal/protocols/entity"

var nodeValidationRules = []func(node entity.Node) error{
	ipCheck,
	pingCheck,
	httpCheck,
	vpnCheck,
}

func GetValidationRules() []func(node entity.Node) error {
	return nodeValidationRules
}
