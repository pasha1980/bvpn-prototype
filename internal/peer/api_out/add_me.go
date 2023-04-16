package api_out

import (
	"bvpn-prototype/internal/peer/api_out/dto"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

func AddMeTo(node entity.Node) {
	method := "/peer/addPeer"
	body, _ := json.Marshal(dto.MakeMeDTO())
	req := fasthttp.AcquireRequest()
	req.SetBody(body)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(node.URL + method)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		time.Sleep(time.Minute)
		AddMeTo(node)
	}

	fasthttp.ReleaseRequest(req)
	if res.StatusCode() != 200 {
		time.Sleep(time.Minute)
		AddMeTo(node)
	}
	fasthttp.ReleaseResponse(res)
}
