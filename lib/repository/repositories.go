package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"tower-defense-api/lib/models"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Repository struct {
	Users interface {
		Create(context.Context, *models.User) error
		GetById(context.Context, int64) (*models.User, error)
	}
	Codes interface {
		Create(context.Context, *models.Code) error
		GetAll(context.Context) ([]*models.Code, error)
	}
	Messages interface {
		Create(context.Context, *models.Message) error
		GetByPlayerId(context.Context, int64) ([]models.Message, error)
		SetRead(context.Context, int64) error
	}
}

func New(db *sql.DB) Repository {
	return Repository{
		Users: &UsersRepository{db},
		Codes: &CodesRepository{db},
		Messages: &MessageRepository{db},
	}
}
