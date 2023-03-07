package internal

import (
	"bvpn-prototype/internal/http/http_in"
	"bvpn-prototype/internal/permatent_tasks"
	"bvpn-prototype/internal/protocols"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/storage/config"
	"fmt"
	"os"
	"strconv"
)

type Kernel struct {
	Address    string
	PublicKey  []byte
	PrivateKey []byte

	URL      string
	HttpPort uint64

	Peers []entity.Node
}

func (k *Kernel) Run() {

	// Save configs
	config.Add(config.Config{
		StorageDirectory: ".",
		URL:              k.URL,
	})

	// Init protocols
	peerProtocol := protocols.GetPeerProtocol()
	chainProtocol := protocols.GetChainProtocol()

	// Check connection to peers
	if _, err := os.Stat("initiate"); err != nil {
		for _, peer := range k.Peers {
			peerProtocol.AddNewPeer(peer)
		}
		os.Create("initiate")
	}

	// Init permanent jobs
	permatent_tasks.Init()

	// Init http controller
	c := http_in.HttpController{
		ChainProtocol: chainProtocol,
		PeerProtocol:  peerProtocol,
	}
	err := http_in.InitHttp(c, ":"+strconv.FormatUint(k.HttpPort, 10), nil)
	if err != nil {
		fmt.Println("Failed to initiate http controller")
		os.Exit(1)
	}
}
