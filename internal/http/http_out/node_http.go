package http_out

import (
	"bvpn-prototype/internal/http/http_dto/chain_stored_dto"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"bytes"
	"encoding/json"
	"net/http"
)

func NewMempoolRecord(stored block_data.ChainStored, nodes ...entity.Node) error {
	var requestDto any
	var method string

	switch stored.Type {
	case block_data.TypeTransaction:
		requestDto = chain_stored_dto.TransactionToDto(stored)
		method = "/addTx"
	case block_data.TypeOffer:
		requestDto = chain_stored_dto.OfferToDto(stored)
		method = "/addOffer"
	case block_data.TypeNodeStatus:
		requestDto = chain_stored_dto.NodeStatusToDto(stored)
		method = "/addTraffic"
	case block_data.TypeTraffic:
		requestDto = chain_stored_dto.TrafficToDto(stored)
		method = "/addNodeStatus"
	}

	body, _ := json.Marshal(requestDto)

	for _, node := range nodes {
		response, err := http.Post(node.URL+method, "application/json", bytes.NewBuffer(body))
		if err != nil {
			// todo
		}

		if response.StatusCode != 200 {
			// todo
		}
	}

	return nil
}

func GetFullChain(node entity.Node) ([]entity.Block, error) {
	return nil, nil
}
