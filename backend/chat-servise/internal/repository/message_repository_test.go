// internal/repository/message_repository_test.go
package repository

import (
	"errors"
	"testing"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)

	tests := []struct {
		name    string
		msg     entity.Message
		mock    func()
		wantErr bool
	}{
		{
			name: "successful message save",
			msg: entity.Message{
				Username: "testuser",
				Message:  "Hello world",
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO chat_messages").
					WithArgs("testuser", "Hello world").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "database error on save",
			msg: entity.Message{
				Username: "testuser",
				Message:  "Hello world",
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO chat_messages").
					WithArgs("testuser", "Hello world").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "empty username",
			msg: entity.Message{
				Username: "",
				Message:  "test",
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO chat_messages").
					WithArgs("", "test").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.SaveMessage(tt.msg)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Message
		wantErr bool
	}{
		{
			name: "successful get messages",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "content"}).
					AddRow(1, "user1", "message 1").
					AddRow(2, "user2", "message 2")
				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
					WillReturnRows(rows)
			},
			want: []entity.Message{
				{ID: 1, Username: "user1", Message: "message 1"},
				{ID: 2, Username: "user2", Message: "message 2"},
			},
			wantErr: false,
		},
		{
			name: "empty result",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "content"})
				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
					WillReturnRows(rows)
			},
			want:    []entity.Message{},
			wantErr: false,
		},
		{
			name: "scan error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "user1")
				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			messages, err := repo.GetMessages()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, messages)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// // internal/repository/message_repository_test.go
// package repository

// import (
// 	"chat-microservice-go/internal/entity"
// 	"errors"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestSaveMessage(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewMessageRepository(db)

// 	tests := []struct {
// 		name    string
// 		msg     entity.Message
// 		mock    func()
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful message save",
// 			msg: entity.Message{
// 				Username: "testuser",
// 				Message:  "Hello world",
// 			},
// 			mock: func() {
// 				mock.ExpectExec("INSERT INTO chat_messages").
// 					WithArgs(32, "testuser", "Hello world").
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "database error on save",
// 			msg: entity.Message{
// 				Username: "testuser",
// 				Message:  "Hello world",
// 			},
// 			mock: func() {
// 				mock.ExpectExec("INSERT INTO chat_messages").
// 					WithArgs(32, "testuser", "Hello world").
// 					WillReturnError(errors.New("database error"))
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "empty message content",
// 			msg: entity.Message{
// 				Username: "user",
// 				Message:  "", // Пустое сообщение
// 			},
// 			mock: func() {
// 				mock.ExpectExec("INSERT INTO chat_messages").
// 					WithArgs(32, "user", "").
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			err := repo.SaveMessage(tt.msg)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }

// func TestGetMessages(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	repo := NewMessageRepository(db)

// 	tests := []struct {
// 		name    string
// 		mock    func()
// 		want    []entity.Message
// 		wantErr bool
// 	}{
// 		{
// 			name: "successful get messages",
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id", "username", "content"}).
// 					AddRow(1, "user1", "message 1").
// 					AddRow(2, "user2", "message 2")
// 				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
// 					WillReturnRows(rows)
// 			},
// 			want: []entity.Message{
// 				{ID: 1, Username: "user1", Message: "message 1"},
// 				{ID: 2, Username: "user2", Message: "message 2"},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "no messages found",
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id", "username", "content"})
// 				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
// 					WillReturnRows(rows)
// 			},
// 			want:    []entity.Message{}, // явно возвращаем пустой слайс вместо nil
// 			wantErr: false,
// 		},
// 		{
// 			name: "database error on get messages",
// 			mock: func() {
// 				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
// 					WillReturnError(errors.New("database error"))
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "scan error",
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id", "username"}). // missing content column
// 											AddRow(1, "user1")
// 				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
// 					WillReturnRows(rows)
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 		{
// 			name: "rows error",
// 			mock: func() {
// 				rows := sqlmock.NewRows([]string{"id", "username", "content"}).
// 					AddRow(1, "user1", "message 1").
// 					RowError(0, errors.New("row error"))
// 				mock.ExpectQuery("SELECT id, username, content FROM chat_messages").
// 					WillReturnRows(rows)
// 			},
// 			want:    nil,
// 			wantErr: true, // убедитесь, что ожидается ошибка
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mock()
// 			got, err := repo.GetMessages()
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.Equal(t, tt.want, got)
// 			assert.NoError(t, mock.ExpectationsWereMet())
// 		})
// 	}
// }
