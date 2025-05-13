package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	user := &domain.User{
		Username:  "testuser",
		Password:  "password",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Password, user.Role, user.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	user := &domain.User{
		Username:  "testuser",
		Password:  "password",
		Role:      "user",
		CreatedAt: time.Now(),
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Password, user.Role, user.CreatedAt).
		WillReturnError(errors.New("database error"))

	id, err := repo.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	username := "testuser"
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role", "created_at"}).
		AddRow(1, username, "password", "user", createdAt)

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, createdAt.Unix(), user.CreatedAt.Unix())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	username := "nonexistent"

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByUsername(context.Background(), username)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	username := "testuser"

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(errors.New("database error"))

	user, err := repo.GetUserByUsername(context.Background(), username)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := int64(1)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "role", "created_at"}).
		AddRow(id, "testuser", "password", "user", createdAt)

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password", user.Password)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, createdAt.Unix(), user.CreatedAt.Unix())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := int64(999)

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(context.Background(), id)
	assert.NoError(t, err) // В вашей реализации возвращается nil, nil при ErrNoRows
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := int64(1)

	mock.ExpectQuery("SELECT id, username, password, role, created_at FROM users WHERE id = \\$1").
		WithArgs(id).
		WillReturnError(errors.New("database error"))

	user, err := repo.GetUserByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}
