package block_data

import "time"

type ConnectionBreak struct {
	Node      string
	Timestamp time.Time
}

const TypeConnectionBreak StoredEntityType = 2
