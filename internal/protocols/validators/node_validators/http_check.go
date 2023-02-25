package node_validators

import (
	"bvpn-prototype/internal/http/http_out"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/protocol_error"
)

func httpCheck(peer entity.Node) error {
	ok := http_out.HealthCheck(peer)
	if !ok {
		return protocol_error.PeerValidationError("HTTP health check failed")
	}

	return nil
}
