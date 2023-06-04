package node_validators

import (
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/peer/domain"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
)

func healthCheck(peer entity.Node) error {
	ok := di.Get("peer_api_gateway").(domain.PeerApiGateway).Peer(peer).HealthCheck()
	if !ok {
		return errors.PeerValidationError("Health check is failed")
	}

	return nil
}
