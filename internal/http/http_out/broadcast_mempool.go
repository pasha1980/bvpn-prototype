package http_out

import (
	"bvpn-prototype/internal/http/http_dto/mempool_data_dto"
	"bvpn-prototype/internal/protocols/entity"
	"bvpn-prototype/internal/protocols/entity/block_data"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func BroadcastMempool(stored block_data.ChainStored, nodes []entity.Node) {
	var requestDto any
	var method string

	switch stored.Type {
	case block_data.TypeTransaction:
		requestDto = mempool_data_dto.TransactionToDto(stored)
		method = "/addTx"
	case block_data.TypeOffer:
		requestDto = mempool_data_dto.OfferToDto(stored)
		method = "/addOffer"
	case block_data.TypeNodeStatus:
		requestDto = mempool_data_dto.NodeStatusToDto(stored)
		method = "/addTraffic"
	case block_data.TypeTraffic:
		requestDto = mempool_data_dto.TrafficToDto(stored)
		method = "/addNodeStatus"
	}

	for _, node := range nodes {
		body, _ := json.Marshal(requestDto)
		req := fasthttp.AcquireRequest()
		req.SetBody(body)
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetRequestURI(node.URL + method)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			time.Sleep(time.Minute)
			BroadcastMempool(stored, []entity.Node{node})
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		if res.StatusCode() != 200 {
			continue
		}
	}
}
