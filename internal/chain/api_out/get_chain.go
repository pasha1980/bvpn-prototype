package api_out

import (
	"bvpn-prototype/internal/http/http_dto"
	"bvpn-prototype/internal/protocols/entity"
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

		var dto http_dto.ChainDto
		err := json.Unmarshal(res.Body(), &dto)
		if err != nil {
			continue
		}

		fasthttp.ReleaseResponse(res)

		for _, blockDto := range dto.Chain {
			chain = append(chain, blockDto.ToEntity())
		}

		if dto.TotalCount < maxBufferSize {
			break
		}

		offset += maxBufferSize
	}

	return chain
}
