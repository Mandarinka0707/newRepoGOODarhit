package entity

import (
	"time"
)

type Message struct {
	ID        int64     `db:"id"`
	TopicID   int64     `db:"topic_id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
