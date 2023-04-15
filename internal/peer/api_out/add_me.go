package api_out

import (
	"bvpn-prototype/internal/http/http_dto"
	"bvpn-prototype/internal/protocols/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func AddMeTo(nodes []entity.Node) {
	method := "/addPeer"
	body, _ := json.Marshal(http_dto.MakeMeDTO())

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
				AddMeTo([]entity.Node{node})
			}()
		}

		fasthttp.ReleaseRequest(req)
		if res.StatusCode() != 200 {
			continue
		}
		fasthttp.ReleaseResponse(res)
	}
}
