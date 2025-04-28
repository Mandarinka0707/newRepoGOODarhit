package entity

import (
	"time"
)

type Session struct {
	UserID    int64     `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
