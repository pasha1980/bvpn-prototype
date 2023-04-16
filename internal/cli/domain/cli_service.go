package domain

import (
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/http"
	"bvpn-prototype/internal/protocol"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"bvpn-prototype/internal/protocol/signer"
	"fmt"
	"os"
	"sort"
	"time"
)

type CliService struct {
}

func (*CliService) Init() {
	if _, err := os.Stat("initiate"); err != nil {
		protocol.InitKeys()

		peerPublic := di.Get("peer_public").(PeerPublicService)
		for _, peer := range config.Get().Peers {
			peerPublic.AddPeer(peer)
		}

		err = di.Get("vpn_public").(VpnPublicService).Init()
		if err != nil {
			panic(err)
		}

		os.Create("initiate")
	}

	chainPublic := di.Get("chain_public").(ChainPublicService)
	go chainPublic.UpdateChain()

	go func() {
		err := http.Init(
			config.Get().HttpPort,
			nil, // todo
		)
		if err != nil {
			panic(err)
		}
	}()
}

func (*CliService) MakeTx(to string, amount float64) {
	chainPublic := di.Get("chain_public").(ChainPublicService)
	utxos, err := chainPublic.GetUTXOs()
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
			tx, _ := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: amount,
				},
			})
			fmt.Println(tx)

			if utxo.Data.(block_data.Transaction).Amount != amount {
				toMe, _ := chainPublic.MakeNew(block_data.ChainStored{
					Type: block_data.TypeTransaction,
					Data: block_data.Transaction{
						From:   utxo.ID.String(),
						To:     signer.GetAddr(),
						Amount: utxo.Data.(block_data.Transaction).Amount - amount,
					},
				})
				fmt.Println(toMe)
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
			tx, _ := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: utxoAmount,
				},
			})
			fmt.Println(tx) // todo
		} else {
			tx, _ := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: sum,
				},
			})
			fmt.Println(tx) // todo

			toMe, _ := chainPublic.MakeNew(block_data.ChainStored{
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

func (*CliService) MakeOffer(price float64) {
	chainPublic := di.Get("chain_public").(ChainPublicService)
	offer, err := chainPublic.MakeNew(block_data.ChainStored{
		Type: block_data.TypeOffer,
		Data: block_data.Offer{
			URL:       config.Get().URL,
			Price:     price,
			Timestamp: time.Now(),
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v\n", offer))
}

func NewCliService() (*CliService, error) {
	return &CliService{}, nil
}
