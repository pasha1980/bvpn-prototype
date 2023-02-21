package block_data

type Transaction struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
	Commission float64 `json:"commission"`
}

const TypeTransaction StoredEntityType = 0
