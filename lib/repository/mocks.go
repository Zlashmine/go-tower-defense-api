package repository

import (
	"context"

	"tower-defense-api/lib/models"
)

func NewMockRepository() Repository {
	return Repository{
		Users:    &MockUsersRepository{},
		Messages: &MockMessagesRepository{},
	}
}

type MockUsersRepository struct{}

type MockMessagesRepository struct{}

func (m *MockMessagesRepository) Create(ctx context.Context, message *models.Message) error {
	return nil
}

func (m *MockMessagesRepository) GetByPlayerId(ctx context.Context, userId int64) ([]models.Message, error) {
	return []models.Message{
		{
			ID:          1,
			UserID:      userId,
			Content:     "test",
			Created:     "2021-06-01T00:00:00Z",
			HasBeenRead: false,
			Sender:      "test",
		},
	}, nil
}

func (m *MockMessagesRepository) SetRead(ctx context.Context, id int64) error {
	return nil
}

func (m *MockUsersRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (m *MockUsersRepository) GetById(ctx context.Context, id int64) (*models.User, error) {
	return &models.User{
		ID:            id,
		Username:      "test",
		AccountStatus: "active",
	}, nil
}
