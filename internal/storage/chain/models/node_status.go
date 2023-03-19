package models

import (
	"bvpn-prototype/internal/protocols/entity/block_data"
	"github.com/google/uuid"
)

type NodeStatus struct {
	ID      uint `gorm:"PRIMARY_KEY"`
	Ref     string
	BlockID uint `gorm:"index"`
	Sign    string
	PubKey  string
	Node    string
	Active  bool
}

func (n NodeStatus) TableName() string {
	return "node_status"
}

func (n *NodeStatus) ToChainStored() block_data.ChainStored {
	id, _ := uuid.FromBytes([]byte(n.Ref))
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeNodeStatus,
		Sign:   n.Sign,
		PubKey: n.PubKey,
		Data: block_data.NodeStatus{
			Node:   n.Node,
			Active: n.Active,
		},
	}
}

func NodeStatusToModel(data block_data.ChainStored, blockId uint) NodeStatus {
	status := data.Data.(block_data.NodeStatus)
	return NodeStatus{
		Ref:     data.ID.String(),
		BlockID: blockId,
		Sign:    data.Sign,
		PubKey:  data.PubKey,
		Node:    status.Node,
		Active:  status.Active,
	}
}
