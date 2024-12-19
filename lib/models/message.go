package models

type Message struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Content     string `json:"content"`
	Created     string `json:"created"`
	HasBeenRead bool   `json:"has_been_read"`
	Sender      string `json:"sender"`
}

type CreateMessagePayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=500"`
	Sender  string `json:"sender" validate:"max=100"`
}
