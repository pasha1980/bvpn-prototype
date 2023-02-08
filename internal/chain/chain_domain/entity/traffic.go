package entity

import (
	"time"
)

type Traffic struct {
	ChainStored

	ID        string    `json:"id"`
	Node      string    `json:"node"`
	Client    string    `json:"client"`
	Gb        float64   `json:"gb"`
	Timestamp time.Time `json:"timestamp"`
}

const TypeTraffic StoredEntityType = 3
