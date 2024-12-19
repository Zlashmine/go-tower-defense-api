package cache

import (
	"context"

	"github.com/stretchr/testify/mock"

	"tower-defense-api/lib/models"
)

func NewMockStore() Store {
	return Store{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct {
	mock.Mock
}

func (m *MockUsersStore) Get(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(id)
	return nil, args.Error(1)
}

func (m *MockUsersStore) Set(ctx context.Context, user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
