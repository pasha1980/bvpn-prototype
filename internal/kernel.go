package internal

import (
	"bvpn-prototype/internal/http/http_in"
	"bvpn-prototype/internal/permatent_tasks"
	"bvpn-prototype/internal/protocols"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/signer"
	"fmt"
	"os"
	"strconv"
)

type Kernel struct {
	URL      string
	HttpPort uint64

	Peers []entity.Node
}

func (k *Kernel) Run() {
	// Init protocols
	peerProtocol := protocols.GetPeerProtocol()
	chainProtocol := protocols.GetChainProtocol()

	// Check if running for the first time
	if _, err := os.Stat("initiate"); err != nil {

		// Add new peers
		for _, peer := range k.Peers {
			peerProtocol.AddNewPeer(peer)
		}

		// Initiate signer package
		signer.Init()

		// Marker
		os.Create("initiate")
	}

	go chainProtocol.UpdateChain()

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

func (k *Kernel) MakeTx(to string, amount float64) {
	protocol := protocols.GetChainProtocol()

	data := protocol.New(block_data.ChainStored{
		Type: block_data.TypeTransaction,
		Data: block_data.Transaction{
			From:   signer.GetAddr(),
			To:     to,
			Amount: amount,
		},
	})

	fmt.Println(data)
}
