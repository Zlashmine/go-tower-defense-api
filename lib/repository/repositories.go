package repository

import (
	"context"
	"database/sql"

	"tower-defense-api/lib/models"
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
}

func New(db *sql.DB) Repository {
	return Repository{
		Users: &UsersRepository{db},
		Codes: &CodesRepository{db},
	}
}
