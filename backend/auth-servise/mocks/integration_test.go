package mocks

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"backend.com/forum/auth-servise/internal/entity"
	"backend.com/forum/auth-servise/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	testDBHost     = "localhost"
	testDBPort     = "5432"
	testDBUser     = "user"
	testDBPassword = "55527"
	testDBName     = "database"
)

var (
	db *sqlx.DB
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testDBHost, testDBPort, testDBUser, testDBPassword, testDBName)

	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	if err := runMigrations(); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}
}

func teardown() {
	if db != nil {
		_, _ = db.Exec("DROP SCHEMA public CASCADE")
		_, _ = db.Exec("CREATE SCHEMA public")
		_ = db.Close()
	}
}

func runMigrations() error {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER NOT NULL,
		token VARCHAR(255) PRIMARY KEY,
		expires_at TIMESTAMP NOT NULL
	);`

	_, err := db.Exec(migrationSQL)
	return err
}

func TestSessionRepositoryIntegration(t *testing.T) {
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	t.Run("Create and Get Session", func(t *testing.T) {

		_, _ = db.Exec("DELETE FROM sessions")

		session := &entity.Session{
			UserID:    1,
			Token:     "test-token-1",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		assert.NoError(t, err)

		retrievedSession, err := repo.GetSessionByToken(ctx, "test-token-1")
		assert.NoError(t, err)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.Token, retrievedSession.Token)
		assert.WithinDuration(t, session.ExpiresAt, retrievedSession.ExpiresAt, time.Second)
	})

	t.Run("Get Non-Existent Session", func(t *testing.T) {

		_, _ = db.Exec("DELETE FROM sessions")

		_, err := repo.GetSessionByToken(ctx, "non-existent-token")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Create Duplicate Session", func(t *testing.T) {

		_, _ = db.Exec("DELETE FROM sessions")

		session := &entity.Session{
			UserID:    1,
			Token:     "duplicate-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		assert.NoError(t, err)

		err = repo.CreateSession(ctx, session)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	})

	t.Run("Multiple Sessions for Same User", func(t *testing.T) {

		_, _ = db.Exec("DELETE FROM sessions")

		session1 := &entity.Session{
			UserID:    1,
			Token:     "token-1",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err := repo.CreateSession(ctx, session1)
		assert.NoError(t, err)

		session2 := &entity.Session{
			UserID:    1,
			Token:     "token-2",
			ExpiresAt: time.Now().Add(48 * time.Hour),
		}
		err = repo.CreateSession(ctx, session2)
		assert.NoError(t, err)

		retrieved1, err := repo.GetSessionByToken(ctx, "token-1")
		assert.NoError(t, err)
		assert.Equal(t, session1.UserID, retrieved1.UserID)

		retrieved2, err := repo.GetSessionByToken(ctx, "token-2")
		assert.NoError(t, err)
		assert.Equal(t, session2.UserID, retrieved2.UserID)
	})
}
