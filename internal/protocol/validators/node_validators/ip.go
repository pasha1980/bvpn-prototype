package node_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/protocol_error"
	"net"
)

func ipCheck(peer entity.Node) error {
	ip := net.ParseIP(peer.IP)
	if ip == nil {
		return protocol_error.PeerValidationError("IP is invalid")
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsInterfaceLocalMulticast() {
		return protocol_error.PeerValidationError("IP range is invalid") // todo: test
	}

	return nil
}
