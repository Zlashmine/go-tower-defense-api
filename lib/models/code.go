package models

type Code struct {
	ID        int64  `json:"id"`
	Code      string `json:"code"`
	Item      string `json:"item"`
	Created   string `json:"created"`
	IsClaimed bool   `json:"is_claimed"`
}

type CreateCodePayload struct {
	Code string `json:"code" validate:"required,max=100"`
	Item string `json:"item" validate:"required,max=100"`
}
