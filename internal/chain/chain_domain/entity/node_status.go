package entity

type NodeStatus struct {
	ChainStored

	Node   string `json:"node"`
	Active bool   `json:"active"`
}

const TypeNodeStatus StoredEntityType = 2