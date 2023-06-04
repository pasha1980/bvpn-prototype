package node_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/errors"
	"net"
)

func ipCheck(peer entity.Node) error {
	ip := net.ParseIP(peer.IP)
	if ip == nil {
		return errors.PeerValidationError("IP is invalid")
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsInterfaceLocalMulticast() {
		return errors.PeerValidationError("IP range is invalid")
	}

	return nil
}
