package block_data

type NodeStatus struct {
	Node   string `json:"node"`
	Active bool   `json:"active"`
}

const TypeNodeStatus StoredEntityType = 2
