package app

import (
	"context"
	"time"

	"net/http"

	"backend.com/forum/forum-servise/internal/domain"
	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/pkg/logger"

	pb "backend.com/forum/proto"
	"github.com/gorilla/websocket"
)

type ForumService struct {
	categoryRepo    repository.CategoryRepository
	topicRepo       repository.TopicRepository
	messageRepo     repository.MessageRepository
	chatMessageRepo repository.ChatMessageRepository
	authClient      pb.AuthServiceClient // gRPC клиент для Auth Service
	logger          *logger.Logger
	upgrader        websocket.Upgrader       // Для веб-сокетов
	chatClients     map[*websocket.Conn]bool // Мапа для хранения подключенных websocket-клиентов
	chatMessageTTL  time.Duration            // Время жизни сообщений в чате
}

type Config struct { // Конфигурационная структура
	AuthServiceAddress string
	ChatMessageTTL     time.Duration
}

func NewForumService(
	categoryRepo repository.CategoryRepository,
	topicRepo repository.TopicRepository,
	messageRepo repository.MessageRepository,
	chatMessageRepo repository.ChatMessageRepository,
	authClient pb.AuthServiceClient,
	cfg *Config,
	logger *logger.Logger) *ForumService {
	return &ForumService{
		categoryRepo:    categoryRepo,
		topicRepo:       topicRepo,
		messageRepo:     messageRepo,
		chatMessageRepo: chatMessageRepo,
		authClient:      authClient,
		logger:          logger,
		upgrader: websocket.Upgrader{ // Инициализация upgrader
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins (for dev)
			},
		},
		chatClients:    make(map[*websocket.Conn]bool),
		chatMessageTTL: cfg.ChatMessageTTL,
	}
}

// Функции для управления категориями, темами, сообщениями, чатом (заготовка)
// ...

// Метод для добавления нового сообщения в чат (обрабатывается в websocketHandler)
func (s *ForumService) AddChatMessage(ctx context.Context, userID int64, content string) error {
	message := &domain.ChatMessage{
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	_, err := s.chatMessageRepo.CreateChatMessage(ctx, message)
	if err != nil {
		s.logger.Errorw("Failed to create chat message in database", "error", err)
		return err
	}
	// Рассылка сообщения всем подключенным клиентам (реализовать)

	s.broadcastChatMessage(message)

	return nil
}

// BroadcastChatMessage рассылает сообщение всем подключенным клиентам
func (s *ForumService) broadcastChatMessage(message *domain.ChatMessage) {
	for client := range s.chatClients {
		// Отправка сообщения через вебсокет (реализовать)
		err := client.WriteJSON(message)
		if err != nil {
			s.logger.Warnw("Websocket write error", "error", err)
			client.Close() // Закрытие соединения при ошибке
			delete(s.chatClients, client)
		}
	}
}

// Функция для удаления устаревших сообщений (cron job)
// func (s *ForumService) cleanupChatMessages() {
// 	ctx := context.Background()
// 	cutoffTime := time.Now().Add(-s.chatMessageTTL) // Вычисление времени удаления
// 	deletedCount, err := s.chatMessageRepo.DeleteChatMessagesBefore(ctx, cutoffTime)

// 	if err != nil {
// 		s.logger.Errorw("Failed to delete old chat messages", "error", err)
// 		return
// 	}

// 	s.logger.Infow("Deleted old chat messages", "count", deletedCount)
// }

// Функция для обработки подключений к веб-сокет чату
func (s *ForumService) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorw("Failed to upgrade to websocket", "error", err)
		return
	}
	s.chatClients[conn] = true // Добавление нового клиента в мапу

	// Обработка сообщений от клиента (реализовать)
	go s.readMessages(conn)
}

// Функция для чтения сообщений от клиента (реализовать)
func (s *ForumService) readMessages(conn *websocket.Conn) {
	defer func() {
		delete(s.chatClients, conn) // Удаление клиента при отключении
		conn.Close()
	}()

	for {
		var message domain.ChatMessage // Структура для получения сообщений

		err := conn.ReadJSON(&message) // Чтение сообщения в формате JSON
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Warnw("Websocket read error", "error", err)
			}
			break
		}

		// Обработка сообщения, например, вызов AddChatMessage
		s.AddChatMessage(context.Background(), message.UserID, message.Content)
	}
}
