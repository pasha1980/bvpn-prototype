package block_data

type Offer struct {
	Node  string  `json:"node"`
	Price float64 `json:"price"`
}

const TypeOffer StoredEntityType = 1
