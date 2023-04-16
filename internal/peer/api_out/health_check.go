package api_out

import (
	"bvpn-prototype/internal/protocol/entity"
	"github.com/valyala/fasthttp"
	"net/http"
)

func HealthCheck(node entity.Node) bool {
	method := "/peer/healthCheck"

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(node.URL + method)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		return false
	}

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	if res.StatusCode() != 200 {
		return false
	}

	return true
}
