package dto

type PaginationDto struct {
	Limit  *int `json:"limit" xml:"limit" form:"limit" query:"limit" validate:"gte=0"`
	Offset *int `json:"offset" xml:"offset" form:"offset" query:"offset" validate:"gte=0"`
}
