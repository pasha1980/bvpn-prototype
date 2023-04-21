package dto

type ProfileShortDto struct {
	ID string `json:"id" validate:"uuid"`
}
