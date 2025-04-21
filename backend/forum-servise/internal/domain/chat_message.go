package domain

import (
	"time"
)

type ChatMessage struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}
