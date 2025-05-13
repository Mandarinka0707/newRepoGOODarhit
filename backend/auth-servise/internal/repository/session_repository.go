package repository

import (
	"context"

	domain "github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/jmoiron/sqlx"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByToken(ctx context.Context, token string) (*domain.Session, error)
}

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, session.UserID, session.Token, session.ExpiresAt)
	return err
}

func (r *sessionRepository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	query := `SELECT user_id, token, expires_at FROM sessions WHERE token = $1`
	session := &domain.Session{}
	err := r.db.GetContext(ctx, session, query, token)
	return session, err
}
