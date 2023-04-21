package domain

import "bvpn-prototype/internal/protocol/entity/block_data"

type ChainPublicService interface {
	GetMyLastOffer() (*block_data.Offer, error)
	SaveTraffic(traffic block_data.Traffic) error
}
