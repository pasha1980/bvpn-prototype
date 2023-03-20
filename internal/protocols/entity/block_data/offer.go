package block_data

type Offer struct {
	Timestamp int64
	Price     float64
}

const TypeOffer StoredEntityType = 1
