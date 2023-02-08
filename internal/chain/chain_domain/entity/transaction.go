package entity

type Transaction struct {
	ChainStored

	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
	Commission float64 `json:"commission"`
}

const TypeTransaction StoredEntityType = 0
