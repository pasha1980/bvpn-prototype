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
	"sort"
	"strconv"
	"time"
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

	go func() {
		// Init http controller
		c := http_in.HttpController{
			ChainProtocol: chainProtocol,
			PeerProtocol:  peerProtocol,
		}
		err := http_in.InitHttp(c, ":"+strconv.FormatUint(k.HttpPort, 10), nil)
		if err != nil {
			// todo
		}
	}()
}

func (k *Kernel) MakeTx(to string, amount float64) {
	protocol := protocols.GetChainProtocol()

	utxos, err := protocol.GetUTXOs()
	if err != nil {
		// todo
		return
	}

	var current float64
	for _, utxo := range utxos {
		current += utxo.Data.(block_data.Transaction).Amount
	}

	if amount > current {
		// todo
		return
	}

	sort.Slice(utxos, func(i int, j int) bool {
		return utxos[i].Data.(block_data.Transaction).Amount > utxos[j].Data.(block_data.Transaction).Amount
	})

	var validForTx []block_data.ChainStored
	var sum float64
	for _, utxo := range utxos {
		if utxo.Data.(block_data.Transaction).Amount >= amount {
			tx := protocol.New(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: amount,
				},
			})
			fmt.Println(tx) // todo

			if utxo.Data.(block_data.Transaction).Amount != amount {
				toMe := protocol.New(block_data.ChainStored{
					Type: block_data.TypeTransaction,
					Data: block_data.Transaction{
						From:   utxo.ID.String(),
						To:     signer.GetAddr(),
						Amount: utxo.Data.(block_data.Transaction).Amount - amount,
					},
				})
				fmt.Println(toMe) // todo
			}
			break
		}

		sum += utxo.Data.(block_data.Transaction).Amount
		validForTx = append(validForTx, utxo)
		if sum > amount {
			break
		}
	}

	for _, utxo := range validForTx {
		utxoAmount := utxo.Data.(block_data.Transaction).Amount
		if sum > utxoAmount {
			tx := protocol.New(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: utxoAmount,
				},
			})
			fmt.Println(tx) // todo
		} else {
			tx := protocol.New(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: sum,
				},
			})
			fmt.Println(tx) // todo

			toMe := protocol.New(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     signer.GetAddr(),
					Amount: utxoAmount - sum,
				},
			})
			fmt.Println(toMe) // todo
		}
		sum -= utxoAmount
	}
}

func (k *Kernel) MakeOffer(price float64) {
	protocol := protocols.GetChainProtocol()
	offer := protocol.New(block_data.ChainStored{
		Type: block_data.TypeOffer,
		Data: block_data.Offer{
			URL:       k.URL,
			Price:     price,
			Timestamp: time.Now().Unix(),
		},
	})

	fmt.Println(fmt.Sprintf("%+v\n", offer))
}
