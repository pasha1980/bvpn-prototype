package domain

import "bvpn-prototype/internal/protocol/entity/block_data"

type ChainPublicService interface {
	GetUTXOs() ([]block_data.ChainStored, error)
	NewEntity(entity block_data.ChainStored) (*block_data.ChainStored, error)
	UpdateChain()
}
