package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend.com/forum/chat-servise/internal/entity"
	"github.com/jmoiron/sqlx"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *entity.Message) (int64, error)
	GetMessage(ctx context.Context, id int64) (*entity.Message, error)
	// ... другие методы
}

type messageRepository struct {
	db *sqlx.DB
}

func NewMessageRepository(db *sqlx.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *entity.Message) (int64, error) {
	query := `INSERT INTO messages (topic_id, user_id, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, message.TopicID, message.UserID, message.Content, message.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *messageRepository) GetMessage(ctx context.Context, id int64) (*entity.Message, error) {
	query := `SELECT id, topic_id, user_id, content, created_at FROM messages WHERE id = $1`
	message := &entity.Message{}
	err := r.db.GetContext(ctx, message, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return message, nil
}
