package domain

import (
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
)

type ChainApiGateway interface {
	To(nodes []entity.Node) ChainApiGateway
	From(node entity.Node) ChainApiGateway
	BroadcastBlock(block entity.Block)
	BroadcastMempool(stored block_data.ChainStored)
	GetChain() []entity.Block
}
