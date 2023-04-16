package api_out

import (
	"bvpn-prototype/internal/chain/api_out/dto"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
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
		requestDto = dto.TransactionToDto(stored)
		method = "/chain/addTx"
	case block_data.TypeOffer:
		requestDto = dto.OfferToDto(stored)
		method = "/chain/addOffer"
	case block_data.TypeTraffic:
		requestDto = dto.ConnectionBreakToDto(stored)
		method = "/chain/addTraffic"
	case block_data.TypeConnectionBreak:
		requestDto = dto.TrafficToDto(stored)
		method = "/chain/addConnectionBreak"
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
			go func() {
				time.Sleep(time.Minute)
				BroadcastMempool(stored, []entity.Node{node})
			}()
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		if res.StatusCode() != 200 {
			continue
		}
	}
}
