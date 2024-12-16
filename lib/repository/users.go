package repository

import (
	"context"
	"database/sql"

	"tower-defense-api/lib/models"
)

type UsersRepository struct {
	db *sql.DB
}

func (repository *UsersRepository) Create(ctx context.Context, payload *models.User) error {
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, created, account_status`

	err := repository.db.
		QueryRowContext(
			ctx,
			query,
			payload.Username).
		Scan(
			&payload.ID,
			&payload.Created,
			&payload.AccountStatus,
		)

	if err != nil {
		return err
	}

	return nil

}

func (repository *UsersRepository) GetById(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, username, created, account_status FROM users WHERE id = $1`

	var user models.User

	err := repository.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.AccountStatus,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
