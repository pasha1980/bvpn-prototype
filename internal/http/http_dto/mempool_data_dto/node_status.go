package mempool_data_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type NodeStatus struct {
	ID     uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Sign   string    `json:"sign" xml:"sign" form:"sign" query:"sign"`
	PubKey string    `json:"pub" xml:"pub" form:"pub" query:"pub"`
	Node   string    `json:"node" xml:"node" form:"node" query:"node"`
	Active bool      `json:"active" xml:"active" form:"active" query:"active"`
}

func (n *NodeStatus) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:     n.ID,
		Sign:   n.Sign,
		PubKey: n.PubKey,
		Type:   block_data.TypeTransaction,
		Data: block_data.NodeStatus{
			Node:   n.Node,
			Active: n.Active,
		},
	}
}

func NodeStatusToDto(entity block_data.ChainStored) NodeStatus {
	nodeStatus := entity.Data.(block_data.NodeStatus)
	return NodeStatus{
		ID:     entity.ID,
		Sign:   entity.Sign,
		PubKey: entity.PubKey,
		Node:   nodeStatus.Node,
		Active: nodeStatus.Active,
	}
}
