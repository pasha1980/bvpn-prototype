package chain_http

type apiResponse struct {
	Error *string        `json:"error"`
	Data  map[string]any `json:"data"`
}
