package repository

import (
	"context"
	"database/sql"

	"tower-defense-api/lib/models"
)

type CodesRepository struct {
	db *sql.DB
}

func (repository *CodesRepository) Create(ctx context.Context, payload *models.Code) error {
	query := `INSERT INTO codes (code, item) VALUES ($1, $2) RETURNING id, created, is_claimed`

	err := repository.db.QueryRowContext(
		ctx,
		query,
		payload.Code,
		payload.Item,
	).Scan(
		&payload.ID,
		&payload.Created,
		&payload.IsClaimed,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *CodesRepository) GetAll(ctx context.Context) ([]*models.Code, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, code, item, created, is_claimed FROM codes")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var codes []*models.Code

	for rows.Next() {
		var code models.Code
		if err := rows.Scan(
			&code.ID,
			&code.Code,
			&code.Item,
			&code.Created,
			&code.IsClaimed,
		); err != nil {
			return nil, err
		}

		codes = append(codes, &code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return codes, nil

}
