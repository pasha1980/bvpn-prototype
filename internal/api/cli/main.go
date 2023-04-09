package cli

import (
	"bvpn-prototype/internal/http/http_in"
	"bvpn-prototype/internal/permatent_tasks"
	"bvpn-prototype/internal/protocols"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bvpn-prototype/internal/protocols/signer"
	"bvpn-prototype/internal/storage/config"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"
)

type CliApi struct {
	ChainProtocol *protocols.ChainProtocol
	PeerProtocol  *protocols.PeerProtocol
	Config        *config.Config
}

func (a *CliApi) Init(detached bool) {
	// Check if running for the first time
	if _, err := os.Stat("initiate"); err != nil {

		// Add new peers
		for _, peer := range a.Config.Peers {
			a.PeerProtocol.AddNewPeer(peer)
		}

		// Initiate signer package
		signer.Init()

		// Marker
		os.Create("initiate")
	}

	go a.ChainProtocol.UpdateChain()

	// Init permanent jobs
	permatent_tasks.Init()

	// Init http server
	go func() {
		err := http_in.InitHttp(":"+strconv.FormatUint(a.Config.HttpPort, 10), nil)
		if err != nil {
			// todo
		}
	}()

	// run in live mode
	if !detached {
		ctlc := make(chan os.Signal)
		signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		<-ctlc
		close(ctlc)
	}
}

func (*CliApi) MakeTx(to string, amount float64) {
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

func (a *CliApi) MakeOffer(price float64) {
	protocol := protocols.GetChainProtocol()
	offer := protocol.New(block_data.ChainStored{
		Type: block_data.TypeOffer,
		Data: block_data.Offer{
			URL:       a.Config.URL,
			Price:     price,
			Timestamp: time.Now(),
		},
	})

	fmt.Println(fmt.Sprintf("%+v\n", offer))
}
