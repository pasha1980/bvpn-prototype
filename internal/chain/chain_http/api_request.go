package chain_http

type apiRequest struct {
	Method    string         `json:"method"`
	Arguments map[string]any `json:"arguments"`
}
