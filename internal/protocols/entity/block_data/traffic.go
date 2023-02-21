package block_data

import (
	"time"
)

type Traffic struct {
	Node      string    `json:"node"`
	Client    string    `json:"client"`
	Gb        float64   `json:"gb"`
	Timestamp time.Time `json:"timestamp"`
}

const TypeTraffic StoredEntityType = 3
