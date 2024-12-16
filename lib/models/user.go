package models

type User struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Created       string `json:"created"`
	AccountStatus string `json:"account_status"`
}

type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
}
