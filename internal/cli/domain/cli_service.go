package domain

import (
	"bvpn-prototype/internal/cli/errors"
	"bvpn-prototype/internal/common"
	"bvpn-prototype/internal/infrastructure/config"
	"bvpn-prototype/internal/infrastructure/di"
	"bvpn-prototype/internal/infrastructure/http"
	"bvpn-prototype/internal/infrastructure/logger"
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

func (*CliService) Init() error {
	if _, err := os.Stat("initiate"); err != nil {
		protocol.InitKeys()
		_, err = os.Create("initiate")
		if err != nil {
			return common.FilesystemError(err.Error())
		}
	}

	peerPublic := di.Get("peer_public").(PeerPublicService)
	for _, peer := range config.Get().Peers {
		peerPublic.AddPeer(peer)
	}

	err := di.Get("vpn_public").(VpnPublicService).Init()
	if err != nil {
		return err
	}

	chainPublic := di.Get("chain_public").(ChainPublicService)
	go chainPublic.UpdateChain()

	go func() {
		err = http.Init(
			config.Get().HttpPort,
			nil, // todo
		)
		if err != nil {
			logger.LogError(err.Error())
		}
	}()

	return nil
}

func (*CliService) MakeTx(to string, amount float64) error {
	chainPublic := di.Get("chain_public").(ChainPublicService)
	utxos, err := chainPublic.GetUTXOs()
	if err != nil {
		return err
	}

	var current float64
	for _, utxo := range utxos {
		current += utxo.Data.(block_data.Transaction).Amount
	}

	if amount > current {
		return errors.InsufficientBalanceError()
	}

	sort.Slice(utxos, func(i int, j int) bool {
		return utxos[i].Data.(block_data.Transaction).Amount > utxos[j].Data.(block_data.Transaction).Amount
	})

	var validForTx []block_data.ChainStored
	var sum float64
	for _, utxo := range utxos {
		if utxo.Data.(block_data.Transaction).Amount >= amount {
			tx, err := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: amount,
				},
			})
			if err != nil {
				return err
			}
			fmt.Println(tx)

			if utxo.Data.(block_data.Transaction).Amount != amount {
				toMe, err := chainPublic.MakeNew(block_data.ChainStored{
					Type: block_data.TypeTransaction,
					Data: block_data.Transaction{
						From:   utxo.ID.String(),
						To:     signer.GetAddr(),
						Amount: utxo.Data.(block_data.Transaction).Amount - amount,
					},
				})
				if err != nil {
					return err
				}
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
			tx, err := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: utxoAmount,
				},
			})
			if err != nil {
				return err
			}
			fmt.Println(tx) // todo
		} else {
			tx, err := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     to,
					Amount: sum,
				},
			})
			if err != nil {
				return err
			}
			fmt.Println(tx) // todo

			toMe, err := chainPublic.MakeNew(block_data.ChainStored{
				Type: block_data.TypeTransaction,
				Data: block_data.Transaction{
					From:   utxo.ID.String(),
					To:     signer.GetAddr(),
					Amount: utxoAmount - sum,
				},
			})
			if err != nil {
				return err
			}
			fmt.Println(toMe) // todo
		}
		sum -= utxoAmount
	}

	return nil
}

func (*CliService) MakeOffer(price float64) error {
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
		return err
	}

	fmt.Println(fmt.Sprintf("%+v\n", offer))
	return nil
}

func NewCliService() (*CliService, error) {
	return &CliService{}, nil
}
