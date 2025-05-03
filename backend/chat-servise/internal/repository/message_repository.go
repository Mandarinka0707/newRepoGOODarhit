package repository

import (
	"chat-microservice-go/internal/entity"
	"database/sql"
	"log"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) SaveMessage(msg entity.Message) error {
	query := `INSERT INTO chat_messages (user_id, username, content) VALUES ($1, $2, $3)`
	_, err := repo.db.Exec(query, 32, msg.Username, msg.Message)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return err
	}
	return nil
}

func (repo *MessageRepository) GetMessages() ([]entity.Message, error) {
	query := `SELECT id, username, content FROM chat_messages`
	rows, err := repo.db.Query(query)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		err := rows.Scan(&msg.ID, &msg.Username, &msg.Message)
		if err != nil {
			log.Printf("Error scanning message: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
