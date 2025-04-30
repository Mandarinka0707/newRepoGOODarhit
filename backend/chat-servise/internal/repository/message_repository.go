// internal/repository/message_repository.go
package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        int64  `gorm:"primaryKey;autoIncrement:true"`
	UserID    int64  `gorm:"column:user_id"`
	Username  string `gorm:"column:username"`
	Content   string
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type MessageRepository interface {
	Create(ctx context.Context, message Message) (*Message, error)
	GetAll(ctx context.Context) ([]Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message Message) (*Message, error) {
	result := r.db.WithContext(ctx).Table("chat_messages").Create(&message)
	return &message, result.Error
}

func (r *messageRepository) GetAll(ctx context.Context) ([]Message, error) {
	var messages []Message
	result := r.db.WithContext(ctx).Table("chat_messages").Order("created_at asc").Find(&messages)
	return messages, result.Error
}
