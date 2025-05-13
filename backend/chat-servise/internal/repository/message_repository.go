// internal/repository/message_repository.go
package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat/internal/entity"
)

type MessageRepository interface {
	SaveMessage(msg entity.Message) error
	GetMessages() ([]entity.Message, error)
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (repo *messageRepository) SaveMessage(msg entity.Message) error {
	query := `INSERT INTO chat_messages (user_id, username, content) VALUES ($1, $2, $3)`
	_, err := repo.db.Exec(query, 32, msg.Username, msg.Message)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return err
	}
	return nil
}

func (repo *messageRepository) GetMessages() ([]entity.Message, error) {
	rows, err := repo.db.Query("SELECT id, username, content FROM chat_messages")
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		// Исправленный маппинг столбцов
		err := rows.Scan(&msg.ID, &msg.Username, &msg.Message)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// // internal/repository/message_repository.go
// package repository

// import (
// 	"chat-microservice-go/internal/entity"
// 	"database/sql"
// 	"log"
// )

// // MessageRepository defines the interface for message repository operations
// type MessageRepository interface {
// 	SaveMessage(msg entity.Message) error
// 	GetMessages() ([]entity.Message, error)
// }

// // messageRepository is the concrete implementation
// type messageRepository struct {
// 	db *sql.DB
// }

// func NewMessageRepository(db *sql.DB) MessageRepository {
// 	return &messageRepository{db: db}
// }

// func (repo *messageRepository) SaveMessage(msg entity.Message) error {
// 	query := `INSERT INTO chat_messages (user_id, username, content) VALUES ($1, $2, $3)`
// 	_, err := repo.db.Exec(query, 32, msg.Username, msg.Message)
// 	if err != nil {
// 		log.Printf("Error saving message: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func (repo *messageRepository) GetMessages() ([]entity.Message, error) {
// 	query := `SELECT id, username, content FROM chat_messages`
// 	rows, err := repo.db.Query(query)
// 	if err != nil {
// 		log.Printf("Error getting messages: %v", err)
// 		return []entity.Message{}, err
// 	}
// 	defer rows.Close()

// 	var messages []entity.Message
// 	for rows.Next() {
// 		var msg entity.Message
// 		err := rows.Scan(&msg.ID, &msg.Username, &msg.Message)
// 		if err != nil {
// 			log.Printf("Error scanning message: %v", err)
// 			return nil, err
// 		}
// 		messages = append(messages, msg)
// 	}
// 	return messages, nil
// }
