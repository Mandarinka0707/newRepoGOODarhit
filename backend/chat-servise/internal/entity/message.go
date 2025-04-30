// internal/entity/message.go
package entity

import "time"

type Message struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:true"`
	UserID    int64     `gorm:"column:user_id"`
	Username  string    `gorm:"column:username"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}
