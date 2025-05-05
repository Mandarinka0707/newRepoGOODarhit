package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"backend.com/forum/auth-servise/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of DB methods we need
type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	argsMock := m.Called(ctx, query, args)
	return argsMock.Get(0).(*sql.Row)
}

func (m *MockDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	argsMock := m.Called(ctx, dest, query, args)
	return argsMock.Error(0)
}

// MockRow is a mock implementation of sql.Row
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// TestUserRepository_CreateUser_Success tests successful user creation
func TestUserRepository_CreateUser_Success(t *testing.T) {
	mockDB := new(MockDB)
	mockRow := new(MockRow)

	now := time.Now()
	user := &entity.User{
		Username:  "testuser",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: now,
	}

	// Mock expectations
	mockRow.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(0).(*int64) = 1 // Set the returned ID to 1
	}).Return(nil)

	mockDB.On("QueryRowContext",
		mock.Anything, // context
		"INSERT INTO users (username, password, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		[]interface{}{user.Username, user.Password, user.Role, user.CreatedAt},
	).Return(mockRow)

	// Create repository with mock DB
	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	id, err := repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestUserRepository_CreateUser_Failure tests user creation failure
func TestUserRepository_CreateUser_Failure(t *testing.T) {
	mockDB := new(MockDB)
	mockRow := new(MockRow)

	now := time.Now()
	user := &entity.User{
		Username:  "testuser",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: now,
	}

	// Mock expectations
	mockRow.On("Scan", mock.Anything).Return(errors.New("database error"))

	mockDB.On("QueryRowContext",
		mock.Anything,
		"INSERT INTO users (username, password, role, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		[]interface{}{user.Username, user.Password, user.Role, user.CreatedAt},
	).Return(mockRow)

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	id, err := repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestUserRepository_GetUserByUsername_Success tests successful user retrieval by username
func TestUserRepository_GetUserByUsername_Success(t *testing.T) {
	mockDB := new(MockDB)

	expectedUser := &entity.User{
		ID:        1,
		Username:  "testuser",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE username = $1",
		[]interface{}{expectedUser.Username},
	).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*entity.User)
		*dest = *expectedUser
	}).Return(nil)

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByUsername(context.Background(), expectedUser.Username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockDB.AssertExpectations(t)
}

// TestUserRepository_GetUserByUsername_NotFound tests user not found scenario
func TestUserRepository_GetUserByUsername_NotFound(t *testing.T) {
	mockDB := new(MockDB)
	username := "nonexistent"

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE username = $1",
		[]interface{}{username},
	).Return(sql.ErrNoRows)

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByUsername(context.Background(), username)

	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, user)
	mockDB.AssertExpectations(t)
}

// TestUserRepository_GetUserByUsername_Failure tests database failure scenario
func TestUserRepository_GetUserByUsername_Failure(t *testing.T) {
	mockDB := new(MockDB)
	username := "testuser"

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE username = $1",
		[]interface{}{username},
	).Return(errors.New("database error"))

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByUsername(context.Background(), username)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockDB.AssertExpectations(t)
}

// TestUserRepository_GetUserByID_Success tests successful user retrieval by ID
func TestUserRepository_GetUserByID_Success(t *testing.T) {
	mockDB := new(MockDB)

	expectedUser := &entity.User{
		ID:        1,
		Username:  "testuser",
		Password:  "hashedpassword",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE id = $1",
		[]interface{}{expectedUser.ID},
	).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*entity.User)
		*dest = *expectedUser
	}).Return(nil)

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByID(context.Background(), expectedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockDB.AssertExpectations(t)
}

// TestUserRepository_GetUserByID_NotFound tests user not found by ID scenario
func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	mockDB := new(MockDB)
	id := int64(999)

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE id = $1",
		[]interface{}{id},
	).Return(sql.ErrNoRows)

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByID(context.Background(), id)

	assert.NoError(t, err)
	assert.Nil(t, user)
	mockDB.AssertExpectations(t)
}

// TestUserRepository_GetUserByID_Failure tests database failure scenario
func TestUserRepository_GetUserByID_Failure(t *testing.T) {
	mockDB := new(MockDB)
	id := int64(1)

	mockDB.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*entity.User"),
		"SELECT id, username, password, role, created_at FROM users WHERE id = $1",
		[]interface{}{id},
	).Return(errors.New("database error"))

	repo := &userRepository{db: &sqlx.DB{DB: mockDB}}

	user, err := repo.GetUserByID(context.Background(), id)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockDB.AssertExpectations(t)
}
