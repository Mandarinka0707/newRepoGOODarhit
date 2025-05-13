package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := &sessionRepository{db: sqlxDB}

	tests := []struct {
		name        string
		session     *entity.Session
		mock        func()
		expectedErr error
	}{
		{
			name: "Success",
			session: &entity.Session{
				UserID:    1,
				Token:     "test-token",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO sessions \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "test-token", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
		{
			name: "Empty Token",
			session: &entity.Session{
				UserID:    1,
				Token:     "",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			mock: func() {
				mock.ExpectExec(`INSERT INTO sessions \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "", sqlmock.AnyArg()).
					WillReturnError(errors.New("empty token"))
			},
			expectedErr: errors.New("empty token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := r.CreateSession(context.Background(), tt.session)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetSessionByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	r := &sessionRepository{db: sqlxDB}

	expiresAt := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		token       string
		mock        func()
		expected    *entity.Session
		expectedErr error
	}{
		{
			name:  "Success",
			token: "valid-token",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "token", "expires_at"}).
					AddRow(1, "valid-token", expiresAt)
				mock.ExpectQuery(`SELECT user_id, token, expires_at FROM sessions WHERE token = \$1`).
					WithArgs("valid-token").
					WillReturnRows(rows)
			},
			expected: &entity.Session{
				UserID:    1,
				Token:     "valid-token",
				ExpiresAt: expiresAt,
			},
			expectedErr: nil,
		},
		{
			name:  "Not Found",
			token: "invalid-token",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, token, expires_at FROM sessions WHERE token = \$1`).
					WithArgs("invalid-token").
					WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := r.GetSessionByToken(context.Background(), tt.token)

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.UserID, result.UserID)
				assert.Equal(t, tt.expected.Token, result.Token)
				assert.WithinDuration(t, tt.expected.ExpiresAt, result.ExpiresAt, time.Second)
			}
		})
	}
}
