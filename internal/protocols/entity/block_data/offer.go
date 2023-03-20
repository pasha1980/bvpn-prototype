package block_data

type Offer struct {
	URL       string
	Timestamp int64
	Price     float64
}

const TypeOffer StoredEntityType = 1
