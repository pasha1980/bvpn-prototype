package entity

type ChainStored struct {
	Author string           `json:"author"`
	Type   StoredEntityType `json:"type"`
}

type StoredEntityType uint8
