package repository

import (
	"context"
	"database/sql"

	"tower-defense-api/lib/models"
)

type MessageRepository struct {
	db *sql.DB
}

func (repository *MessageRepository) Create(ctx context.Context, payload *models.Message) error {
	query := `
	INSERT INTO messages (user_id, content, sender) VALUES ($1, $2, $3) RETURNING id, created, has_been_read
	`

	err := repository.db.QueryRowContext(
		ctx,
		query,
		payload.UserID,
		payload.Content,
		payload.Sender,
	).Scan(
		&payload.ID,
		&payload.Created,
		&payload.HasBeenRead,
	)

	return err
}

func (repository *MessageRepository) GetByPlayerId(ctx context.Context, playerId int64) ([]models.Message, error) {
	query := `
	SELECT m.id, m.user_id, m.content, m.created, m.has_been_read, m.sender FROM messages m
	LEFT JOIN users u ON m.user_id = u.id 
	WHERE m.user_id = $1 AND m.has_been_read = false
	ORDER BY m.created DESC
	`

	rows, err := repository.db.QueryContext(ctx, query, playerId)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	messages := []models.Message{}

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.Created,
			&message.HasBeenRead,
			&message.Sender,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (repository *MessageRepository) SetRead(ctx context.Context, id int64) error {
	query := `UPDATE messages SET has_been_read = true WHERE id = $1`

	_, err := repository.db.ExecContext(ctx, query, id)

	return err
}
