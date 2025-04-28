package usecase

import (
	"context"
	"net/http"
	"time"

	"backend.com/forum/chat-servise/internal/entity"
	"backend.com/forum/chat-servise/internal/repository"
	"backend.com/forum/chat-servise/pkg/logger"
	"github.com/gorilla/websocket"
)

type ChatUsecase struct {
	chatMessageRepo repository.ChatMessageRepository
	logger          *logger.Logger
	upgrader        websocket.Upgrader
	chatClients     map[*websocket.Conn]bool
	chatMessageTTL  time.Duration
}

type ChatConfig struct {
	MessageTTL time.Duration
}

func NewChatUsecase(
	chatMessageRepo repository.ChatMessageRepository,
	cfg *ChatConfig,
	logger *logger.Logger,
) *ChatUsecase {
	return &ChatUsecase{
		chatMessageRepo: chatMessageRepo,
		logger:          logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		chatClients:    make(map[*websocket.Conn]bool),
		chatMessageTTL: cfg.MessageTTL,
	}
}

func (uc *ChatUsecase) AddChatMessage(ctx context.Context, userID int64, content string) error {
	message := &entity.ChatMessage{
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	_, err := uc.chatMessageRepo.CreateChatMessage(ctx, message)
	if err != nil {
		uc.logger.Errorw("Failed to create chat message", "error", err)
		return err
	}
	uc.broadcastChatMessage(message)
	return nil
}

func (uc *ChatUsecase) broadcastChatMessage(message *entity.ChatMessage) {
	for client := range uc.chatClients {
		if err := client.WriteJSON(message); err != nil {
			uc.logger.Warnw("Websocket write error", "error", err)
			_ = client.Close()
			delete(uc.chatClients, client)
		}
	}
}

// Добавляем метод для доступа к upgrader
func (uc *ChatUsecase) Upgrader() *websocket.Upgrader {
	return &uc.upgrader
}
func (uc *ChatUsecase) HandleWebSocket(conn *websocket.Conn) {
	uc.chatClients[conn] = true
	defer func() {
		delete(uc.chatClients, conn)
		_ = conn.Close()
	}()

	for {
		var msg entity.ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				uc.logger.Warnw("Websocket error", "error", err)
			}
			break
		}

		if err := uc.AddChatMessage(context.Background(), msg.UserID, msg.Content); err != nil {
			uc.logger.Warnw("Failed to process chat message", "error", err)
		}
	}
}
