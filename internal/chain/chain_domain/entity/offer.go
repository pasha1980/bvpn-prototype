package entity

type Offer struct {
	ChainStored

	Node string  `json:"node"`
	Price    float64 `json:"price"`
}

const TypeOffer StoredEntityType = 1
