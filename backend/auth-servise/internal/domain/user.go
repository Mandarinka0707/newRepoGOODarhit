package domain

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"` // Хранится хеш пароля
	Role      string    `db:"role"`     // "admin", "user"
	CreatedAt time.Time `db:"created_at"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
