package infrastructure

import (
	"bvpn-prototype/tests/testing"
	"github.com/go-ping/ping"
	"os/exec"
	"time"
)

func It_is_possible_to_use_ip_command(t *testing.T) error {
	_, err := exec.LookPath("ip")
	t.Assert(err == nil, "Unable to use `ip` command")
	return nil
}

func It_is_stable_internet_connection(t *testing.T) error {
	pingUrl := "google.com"

	pinger, _ := ping.NewPinger(pingUrl)
	pinger.Timeout = time.Second
	pinger.Run()
	t.Assert(pinger.Statistics().PacketLoss == 0.0, "Unstable internet connection")

	return nil
}
