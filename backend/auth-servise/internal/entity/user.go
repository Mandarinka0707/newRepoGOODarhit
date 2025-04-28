package entity

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
