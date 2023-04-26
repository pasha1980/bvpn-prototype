package models

import (
	"bvpn-prototype/internal/protocol/entity/block_data"
	"encoding/json"
	"github.com/google/uuid"
)

type UndefinedData struct {
	ID      uint `gorm:"PRIMARY_KEY"`
	Ref     string
	BlockID uint `gorm:"index"`
	Sign    string
	PubKey  string
	Fields  string
}

func (d UndefinedData) TableName() string {
	return "undefined_data"
}

func (d *UndefinedData) ToChainStored() block_data.ChainStored {
	id, _ := uuid.Parse(d.Ref)

	var fields map[string]any
	json.Unmarshal([]byte(d.Fields), &fields)
	return block_data.ChainStored{
		ID:     id,
		Type:   block_data.TypeTraffic,
		Sign:   d.Sign,
		PubKey: d.PubKey,
		Data:   fields,
	}
}

func UndefinedDataModel(data block_data.ChainStored, blockId uint) UndefinedData {
	dataMap := data.Data.(map[string]any)
	encoded, _ := json.Marshal(dataMap)
	return UndefinedData{
		Ref:     data.ID.String(),
		BlockID: blockId,
		Sign:    data.Sign,
		PubKey:  data.PubKey,
		Fields:  string(encoded),
	}
}
