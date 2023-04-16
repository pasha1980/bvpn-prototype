package dto

type ChainDto struct {
	Chain      []BlockDto `json:"chain"`
	TotalCount int        `json:"totalCount"`
}
