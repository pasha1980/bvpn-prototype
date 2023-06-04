package api_gateway

import (
	"bvpn-prototype/internal/peer/api_gateway/dto"
	"bvpn-prototype/internal/peer/domain"
	"bvpn-prototype/internal/protocol/entity"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type PeerApiGatewayImpl struct {
	peer entity.Node
}

func (g *PeerApiGatewayImpl) Peer(node entity.Node) domain.PeerApiGateway {
	g.peer = node
	return g
}

func (g *PeerApiGatewayImpl) AddMe() {
	method := "/peer/addPeer"
	body, _ := json.Marshal(dto.MakeMeDTO())
	req := fasthttp.AcquireRequest()
	req.SetBody(body)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(g.peer.URL + method)
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		time.Sleep(time.Minute)
		g.AddMe()
	}

	fasthttp.ReleaseRequest(req)
	if res.StatusCode() != 200 {
		time.Sleep(time.Minute)
		g.AddMe()
	}
	fasthttp.ReleaseResponse(res)
}

func (g *PeerApiGatewayImpl) HealthCheck() bool {
	method := "/peer/healthCheck"

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(g.peer.URL + method)
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

func NewPeerApiGateway() (*PeerApiGatewayImpl, error) {
	return &PeerApiGatewayImpl{}, nil
}
