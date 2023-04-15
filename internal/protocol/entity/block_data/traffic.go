package block_data

import (
	"time"
)

type Traffic struct {
	Node      string
	Client    string
	Gb        float64
	Timestamp time.Time
}

const TypeTraffic StoredEntityType = 3
