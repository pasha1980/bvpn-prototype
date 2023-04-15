package block_data

import "time"

type Offer struct {
	URL       string
	Timestamp time.Time
	Price     float64
}

const TypeOffer StoredEntityType = 1
