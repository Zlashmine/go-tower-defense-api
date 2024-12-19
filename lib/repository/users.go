package repository

import (
	"context"
	"database/sql"

	"tower-defense-api/lib/models"
)

type UsersRepository struct {
	Db *sql.DB
}

func (repository *UsersRepository) Create(ctx context.Context, payload *models.User) error {
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, created, account_status`

	err := repository.Db.
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User

	err := repository.Db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.AccountStatus,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
