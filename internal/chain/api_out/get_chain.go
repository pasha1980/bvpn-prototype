package api_out

import (
	"bvpn-prototype/internal/chain/api_out/dto"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func GetFullChain(node entity.Node) []entity.Block {
	var chain []entity.Block

	const maxBufferSize = 50
	const method = "/getChain"

	offset := 0
	limit := maxBufferSize

	for {
		body, _ := json.Marshal(map[string]int{
			"offset": offset,
			"limit":  limit,
		})

		req := fasthttp.AcquireRequest()
		req.SetBody(body)
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetRequestURI(node.URL + method)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			continue
		}

		fasthttp.ReleaseRequest(req)
		if res.StatusCode() != 200 {
			continue
		}

		var d dto.ChainDto
		err := json.Unmarshal(res.Body(), &d)
		if err != nil {
			continue
		}

		fasthttp.ReleaseResponse(res)

		for _, blockDto := range d.Chain {
			chain = append(chain, blockDto.ToEntity())
		}

		if d.TotalCount < maxBufferSize {
			break
		}

		offset += maxBufferSize
	}

	return chain
}
