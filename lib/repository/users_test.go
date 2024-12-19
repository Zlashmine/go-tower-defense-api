package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"tower-defense-api/lib/models"
	"tower-defense-api/lib/repository"
)

var (
	mockDB               *sql.DB
	sqlMock              sqlmock.Sqlmock
	mockCtx              context.Context
	QueryTimeoutDuration = 10 * time.Second // Example timeout
)

func setupTests(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockDB = db
	sqlMock = mock
	mockCtx = context.Background()
}

func TestUsersRepository_Create(t *testing.T) {
	setupTests(t)
	defer mockDB.Close()

	user := &models.User{Username: "test_user"}

	sqlMock.ExpectQuery("INSERT INTO users (.+) VALUES (.+) RETURNING id, created, account_status").
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created", "account_status"}).AddRow(1, time.Now(), "active"))

	userRepository := repository.UsersRepository{Db: mockDB}
	err := userRepository.Create(mockCtx, user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotNil(t, user.Created)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUsersRepository_Create_Error(t *testing.T) {
	setupTests(t)
	defer mockDB.Close()

	user := &models.User{Username: "test_user"}

	sqlMock.ExpectQuery("INSERT INTO users (.+) VALUES (.+) RETURNING id, created, account_status").
		WithArgs(user.Username).
		WillReturnError(errors.New("database error"))

	userRepository := repository.UsersRepository{Db: mockDB}
	err := userRepository.Create(mockCtx, user)

	assert.Error(t, err)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUsersRepository_GetById_Found(t *testing.T) {
	setupTests(t)
	defer mockDB.Close()

	expectedID := int64(1)
	expectedUsername := "test_user"

	rows := sqlmock.NewRows([]string{"id", "username", "created", "account_status"}).
		AddRow(expectedID, expectedUsername, time.Now(), "active")

	sqlMock.ExpectQuery("SELECT id, username, created, account_status FROM users WHERE id = ?").
		WithArgs(expectedID).
		WillReturnRows(rows)

	// TODO Use Mock Repository and inject mockDB
	userRepository := repository.UsersRepository{Db: mockDB}
	user, err := userRepository.GetById(mockCtx, expectedID)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, user.ID)
	assert.Equal(t, expectedUsername, user.Username)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUsersRepository_GetById_NotFound(t *testing.T) {
	setupTests(t)
	defer mockDB.Close()

	expectedID := int64(1)

	sqlMock.ExpectQuery("SELECT id, username, created, account_status FROM users WHERE id = ?").
		WithArgs(expectedID).
		WillReturnError(sql.ErrNoRows)

	userRepository := repository.UsersRepository{Db: mockDB}
	user, err := userRepository.GetById(mockCtx, expectedID)

	assert.Equal(t, repository.ErrNotFound, err)
	assert.Nil(t, user)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUsersRepository_GetById_Error(t *testing.T) {
	setupTests(t)
	defer mockDB.Close()

	expectedID := int64(1)

	sqlMock.ExpectQuery("SELECT id, username, created, account_status FROM users WHERE id = ?").
		WithArgs(expectedID).
		WillReturnError(errors.New("some error"))

	userRepository := repository.UsersRepository{Db: mockDB}
	user, err := userRepository.GetById(mockCtx, expectedID)

	assert.Error(t, err)
	assert.Nil(t, user)

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
