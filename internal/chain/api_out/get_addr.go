package api_out

import (
	"bvpn-prototype/internal/protocols/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func GetAddr(node entity.Node) (string, error) {
	const method = "/getAddress"
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(node.URL + method)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return "", err
	}

	fasthttp.ReleaseRequest(req)
	if res.StatusCode() != 200 {
		return "", nil // todo
	}
	fasthttp.ReleaseResponse(res)

	var response map[string]string
	err := json.Unmarshal(res.Body(), &response)
	if err != nil {
		return "", err
	}

	addr, ok := response["address"]
	if !ok {
		return "", nil // todo
	}

	return addr, nil
}
