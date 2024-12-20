package notifications

import (
	"embed"

	"tower-defense-api/lib/models"
)

const (
	FromName            = "Maze Defenders"
	maxRetires          = 3
	MessageTemplate = "message.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(message *models.Message, isSandbox bool) (int, error)
}
