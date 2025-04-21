package repository

import (
	"context"
	"time"

	"backend.com/forum/forum-servise/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ChatMessageRepository interface {
	CreateChatMessage(ctx context.Context, message *domain.ChatMessage) (int64, error)
	DeleteChatMessagesBefore(ctx context.Context, cutoffTime time.Time) (int64, error) // Add this method
}

type chatMessageRepository struct {
	db *sqlx.DB
}

func NewChatMessageRepository(db *sqlx.DB) ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) CreateChatMessage(ctx context.Context, message *domain.ChatMessage) (int64, error) {
	query := `INSERT INTO chat_messages (user_id, content, created_at) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, message.UserID, message.Content, message.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *chatMessageRepository) DeleteChatMessagesBefore(ctx context.Context, cutoffTime time.Time) (int64, error) {
	query := `DELETE FROM chat_messages WHERE created_at < $1`
	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
