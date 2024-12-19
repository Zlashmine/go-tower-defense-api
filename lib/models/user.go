package models

type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	Created       string    `json:"created"`
	AccountStatus string    `json:"account_status"`
	Messages      []Message `json:"messages"`
}

// CreateUserPayload godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Tags			users
type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=2,max=100"`
}
