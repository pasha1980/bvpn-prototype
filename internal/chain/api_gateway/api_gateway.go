package api_gateway

import (
	"bvpn-prototype/internal/chain/api_gateway/dto"
	"bvpn-prototype/internal/chain/domain"
	"bvpn-prototype/internal/protocol/entity"
	"bvpn-prototype/internal/protocol/entity/block_data"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type ChainApiGatewayImpl struct {
	toNodes  []entity.Node
	fromNode entity.Node
}

func (g *ChainApiGatewayImpl) To(nodes []entity.Node) domain.ChainApiGateway {
	g.toNodes = nodes
	return g
}

func (g *ChainApiGatewayImpl) From(node entity.Node) domain.ChainApiGateway {
	g.fromNode = node
	return g
}

func (g *ChainApiGatewayImpl) BroadcastBlock(block entity.Block) {
	method := "/chain/addBlock"
	body, _ := json.Marshal(dto.BlockToDto(block))

	for _, node := range g.toNodes {
		req := fasthttp.AcquireRequest()
		req.SetBody(body)
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetRequestURI(node.URL + method)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			go func() {
				time.Sleep(time.Minute)
				g.BroadcastBlock(block)
			}()
		}

		fasthttp.ReleaseRequest(req)
		if res.StatusCode() != 200 {
			continue
		}
		fasthttp.ReleaseResponse(res)
	}

	g.toNodes = []entity.Node{}
}

func (g *ChainApiGatewayImpl) BroadcastMempool(stored block_data.ChainStored) {
	var requestDto any
	var method string

	switch stored.Type {
	case block_data.TypeTransaction:
		requestDto = dto.TransactionToDto(stored)
		method = "/chain/addTx"
	case block_data.TypeOffer:
		requestDto = dto.OfferToDto(stored)
		method = "/chain/addOffer"
	case block_data.TypeTraffic:
		requestDto = dto.ConnectionBreakToDto(stored)
		method = "/chain/addTraffic"
	case block_data.TypeConnectionBreak:
		requestDto = dto.TrafficToDto(stored)
		method = "/chain/addConnectionBreak"
	default:
		return
	}

	for _, node := range g.toNodes {
		body, _ := json.Marshal(requestDto)
		req := fasthttp.AcquireRequest()
		req.SetBody(body)
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetRequestURI(node.URL + method)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			go func() {
				time.Sleep(time.Minute)
				g.BroadcastMempool(stored)
			}()
		}
		fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		if res.StatusCode() != 200 {
			continue
		}
	}

	g.toNodes = []entity.Node{}
}

func (g *ChainApiGatewayImpl) GetChain() []entity.Block {
	var chain []entity.Block

	const maxBufferSize = 50
	const method = "/chain/getChain"

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
		req.SetRequestURI(g.fromNode.URL + method)
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

func NewChainApiGateway() (*ChainApiGatewayImpl, error) {
	return &ChainApiGatewayImpl{}, nil
}
