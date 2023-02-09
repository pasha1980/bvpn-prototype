package chain_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Peer struct {
	IP  string
	URL string
}

func (p *Peer) Do(method string, arguments map[string]any) (map[string]any, error) {
	body := apiRequest{
		Method:    method,
		Arguments: arguments,
	}
	j, _ := json.Marshal(body)

	reader := bytes.NewReader(j)
	r, _ := http.NewRequest(http.MethodPost, p.URL, reader)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var decodedResponse apiResponse
	err = json.NewDecoder(r.Body).Decode(&decodedResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest || decodedResponse.Error != nil {
		return nil, errors.New(*decodedResponse.Error)
	}

	return decodedResponse.Data, nil
}
