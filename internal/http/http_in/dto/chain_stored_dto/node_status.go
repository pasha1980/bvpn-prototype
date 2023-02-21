package chain_stored_dto

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type NodeStatus struct {
	ID     uuid.UUID `json:"id" xml:"id" form:"id" query:"id"`
	Node   string    `json:"node" xml:"node" form:"node" query:"node"`
	Active bool      `json:"active" xml:"active" form:"active" query:"active"`
}

func (n *NodeStatus) ToEntity() block_data.ChainStored {
	return block_data.ChainStored{
		ID:   n.ID,
		Type: block_data.TypeTransaction,
		Data: block_data.NodeStatus{
			Node:   n.Node,
			Active: n.Active,
		},
	}
}
