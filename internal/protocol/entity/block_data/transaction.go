package block_data

type Transaction struct {
	From   string
	To     string
	Amount float64
}

const TypeTransaction StoredEntityType = 0
