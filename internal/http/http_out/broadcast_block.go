package http_out

import (
	"bvpn-prototype/internal/http/http_dto"
	"bvpn-prototype/internal/protocols/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func BroadcastBlock(block entity.Block, nodes []entity.Node) { // todo: In method
	method := "/addBlock"
	body, _ := json.Marshal(http_dto.BlockToDto(block))

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
