package repository

import (
	"context"
	"database/sql"
	"errors"

	"backend.com/forum/forum-servise/internal/entity"
	"github.com/jmoiron/sqlx"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *entity.Topic) (int64, error)
	GetTopic(ctx context.Context, id int64) (*entity.Topic, error)
}

type topicRepository struct {
	db *sqlx.DB
}

func NewTopicRepository(db *sqlx.DB) TopicRepository {
	return &topicRepository{db: db}
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *entity.Topic) (int64, error) {
	query := `INSERT INTO topics (category_id, title, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, topic.CategoryID, topic.Title, topic.UserID, topic.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *topicRepository) GetTopic(ctx context.Context, id int64) (*entity.Topic, error) {
	query := `SELECT id, category_id, title, user_id, created_at FROM topics WHERE id = $1`
	topic := &entity.Topic{}
	err := r.db.GetContext(ctx, topic, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return topic, nil
}
