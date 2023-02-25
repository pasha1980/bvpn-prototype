package http_out

import (
	"bvpn-prototype/internal/protocols/entity"
	"github.com/valyala/fasthttp"
	"net/http"
)

func HealthCheck(node entity.Node) bool { // todo: In method
	method := "/healthCheck"

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
