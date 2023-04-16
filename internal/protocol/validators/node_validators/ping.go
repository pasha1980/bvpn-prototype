package node_validators

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/protocol_error"
	"github.com/go-ping/ping"
	"time"
)

func pingCheck(peer entity.Node) error {
	pingUrl := peer.URL
	if pingUrl == "" {
		pingUrl = peer.IP
	}

	pinger, _ := ping.NewPinger(pingUrl)
	pinger.Timeout = time.Second
	pinger.Run()

	stat := pinger.Statistics()
	if stat.PacketsRecv == 0 {
		return protocol_error.PeerValidationError("Peer unavailable")
	}

	return nil
}
