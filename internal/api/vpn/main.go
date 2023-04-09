package vpn

import "bvpn-prototype/internal/protocols"

type VpnApi struct {
	ChainProtocol *protocols.ChainProtocol
	PeerProtocol  *protocols.PeerProtocol
}
