package cache

import (
	"context"

	"tower-defense-api/lib/models"
)

type Store struct {
	Users interface {
		Get(context.Context, int64) (*models.User, error)
		Set(context.Context, *models.User) error
	}
}
