package api_out

import (
	"bvpn-prototype/internal/chain/api_out/dto"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func BroadcastBlock(block entity.Block, nodes []entity.Node) {
	method := "/chain/addBlock"
	body, _ := json.Marshal(dto.BlockToDto(block))

	for _, node := range nodes {
		req := fasthttp.AcquireRequest()
		req.SetBody(body)
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetRequestURI(node.URL + method)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			go func() {
				time.Sleep(time.Minute)
				BroadcastBlock(block, []entity.Node{node})
			}()
		}

		fasthttp.ReleaseRequest(req)
		if res.StatusCode() != 200 {
			continue
		}
		fasthttp.ReleaseResponse(res)
	}

}
