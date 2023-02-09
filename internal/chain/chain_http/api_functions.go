package chain_http

var apiFunctions = map[string]func(arguments map[string]any) (map[string]any, error){}
