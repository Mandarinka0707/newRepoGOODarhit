package app

import (
	"context"
	"strings"
	"time"

	"net/http"

	"backend.com/forum/forum-servise/internal/domain"
	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/pkg/logger"

	pb "backend.com/forum/proto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ForumService struct {
	categoryRepo    repository.CategoryRepository
	topicRepo       repository.TopicRepository
	messageRepo     repository.MessageRepository
	chatMessageRepo repository.ChatMessageRepository
	postRepo        repository.PostRepository
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
	postRepo repository.PostRepository,
	authClient pb.AuthServiceClient,
	cfg *Config,
	logger *logger.Logger) *ForumService {
	return &ForumService{
		categoryRepo:    categoryRepo,
		topicRepo:       topicRepo,
		messageRepo:     messageRepo,
		chatMessageRepo: chatMessageRepo,
		postRepo:        postRepo,
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

// Метод для добавления нового сообщения в чат (обрабатывается в websocketHandle
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
	// Рассылка сообщения всем подключенным клиентам
	s.broadcastChatMessage(message)

	return nil
}

// BroadcastChatMessage рассылает сообщение всем подключенным клиентам
func (s *ForumService) broadcastChatMessage(message *domain.ChatMessage) {
	for client := range s.chatClients {
		// Отправка сообщения через вебсокет
		err := client.WriteJSON(message)
		if err != nil {
			s.logger.Warnw("Websocket write error", "error", err)
			client.Close() // Закрытие соединения при ошибке
			delete(s.chatClients, client)
		}
	}
}

// Функция для обработки подключений к веб-сокет чат
func (s *ForumService) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorw("Failed to upgrade to websocket", "error", err)
		return
	}
	s.chatClients[conn] = true // Добавление нового клиента в мапу

	// Обработка сообщений от клиента
	go s.readMessages(conn)
}

// Функция для чтения сообщений от клиента
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
func (s *ForumService) CreatePost(c *gin.Context) {
	// 1. Check authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// 2. Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 3. Validate token and get user info
	authResponse, err := s.authClient.ValidateToken(c.Request.Context(), &pb.ValidateTokenRequest{
		Token: tokenString,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 4. Verify the user exists in the database
	userResponse, err := s.authClient.GetUser(c.Request.Context(), &pb.GetUserRequest{
		Id: authResponse.UserId,
	})
	if err != nil || userResponse.User == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// 5. Parse input
	var post struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// 6. Create post in DB
	newPost := &domain.Post{
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  authResponse.UserId,
		CreatedAt: time.Now(),
	}

	id, err := s.postRepo.CreatePost(c.Request.Context(), newPost)
	if err != nil {
		s.logger.Errorw("Failed to create post",
			"error", err,
			"post", newPost,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create post",
			"details": err.Error(),
		})
		return
	}

	// 7. Success response
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Post created successfully",
		"post": gin.H{
			"title":     newPost.Title,
			"content":   newPost.Content,
			"author_id": newPost.AuthorID,
		},
	})
}

// Вспомогательная функция для валидации токена
// func (s *ForumService) validateToken(tokenString string) (int64, error) {
// 	// Реализуйте проверку токена через ваш authClient
// 	// Это пример - замените на вашу реальную реализацию
// 	return 123, nil // Возвращаем ID пользователя
// }

// GetPosts возвращает список всех постов
func (s *ForumService) GetPosts(c *gin.Context) {
	posts, err := s.postRepo.GetPosts(c.Request.Context())
	if err != nil {
		s.logger.Errorw("Failed to get posts", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	// Validate posts before returning
	for i := range posts {
		if posts[i].ID == 0 {
			s.logger.Warnw("Post with missing ID found", "post", posts[i])
			posts[i].ID = -1 // Or generate a temporary ID
		}
	}

	c.JSON(http.StatusOK, posts)
}

// func NewForumService(
// 	categoryRepo repository.CategoryRepository,
// 	topicRepo repository.TopicRepository,
// 	messageRepo repository.MessageRepository,
// 	chatMessageRepo repository.ChatMessageRepository,
// 	authClient pb.AuthServiceClient,
// 	cfg *Config,
// 	logger *logger.Logger) *ForumService {
// 	return &ForumService{
// 		categoryRepo:    categoryRepo,
// 		topicRepo:       topicRepo,
// 		messageRepo:     messageRepo,
// 		chatMessageRepo: chatMessageRepo,
// 		authClient:      authClient,
// 		logger:          logger,
// 		upgrader: websocket.Upgrader{ // Инициализация upgrader
// 			CheckOrigin: func(r *http.Request) bool {
// 				return true // Allow all origins (for dev)
// 			},
// 		},
// 		chatClients:    make(map[*websocket.Conn]bool),
// 		chatMessageTTL: cfg.ChatMessageTTL,
// 	}
// }

// // Функции для управления категориями, темами, сообщениями, чатом (заготовка)
// // ...

// // Метод для добавления нового сообщения в чат (обрабатывается в websocketHandler)
// func (s *ForumService) AddChatMessage(ctx context.Context, userID int64, content string) error {
// 	message := &domain.ChatMessage{
// 		UserID:    userID,
// 		Content:   content,
// 		CreatedAt: time.Now(),
// 	}

// 	_, err := s.chatMessageRepo.CreateChatMessage(ctx, message)
// 	if err != nil {
// 		s.logger.Errorw("Failed to create chat message in database", "error", err)
// 		return err
// 	}
// 	// Рассылка сообщения всем подключенным клиентам (реализовать)

// 	s.broadcastChatMessage(message)

// 	return nil
// }

// // BroadcastChatMessage рассылает сообщение всем подключенным клиентам
// func (s *ForumService) broadcastChatMessage(message *domain.ChatMessage) {
// 	for client := range s.chatClients {
// 		// Отправка сообщения через вебсокет (реализовать)
// 		err := client.WriteJSON(message)
// 		if err != nil {
// 			s.logger.Warnw("Websocket write error", "error", err)
// 			client.Close() // Закрытие соединения при ошибке
// 			delete(s.chatClients, client)
// 		}
// 	}
// }

// // Функция для удаления устаревших сообщений (cron job)
// // func (s *ForumService) cleanupChatMessages() {
// // 	ctx := context.Background()
// // 	cutoffTime := time.Now().Add(-s.chatMessageTTL) // Вычисление времени удаления
// // 	deletedCount, err := s.chatMessageRepo.DeleteChatMessagesBefore(ctx, cutoffTime)

// // 	if err != nil {
// // 		s.logger.Errorw("Failed to delete old chat messages", "error", err)
// // 		return
// // 	}

// // 	s.logger.Infow("Deleted old chat messages", "count", deletedCount)
// // }

// // Функция для обработки подключений к веб-сокет чату
// func (s *ForumService) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := s.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		s.logger.Errorw("Failed to upgrade to websocket", "error", err)
// 		return
// 	}
// 	s.chatClients[conn] = true // Добавление нового клиента в мапу

// 	// Обработка сообщений от клиента (реализовать)
// 	go s.readMessages(conn)
// }

// // Функция для чтения сообщений от клиента (реализовать)
// func (s *ForumService) readMessages(conn *websocket.Conn) {
// 	defer func() {
// 		delete(s.chatClients, conn) // Удаление клиента при отключении
// 		conn.Close()
// 	}()

// 	for {
// 		var message domain.ChatMessage // Структура для получения сообщений

// 		err := conn.ReadJSON(&message) // Чтение сообщения в формате JSON
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				s.logger.Warnw("Websocket read error", "error", err)
// 			}
// 			break
// 		}

// 		// Обработка сообщения, например, вызов AddChatMessage
// 		s.AddChatMessage(context.Background(), message.UserID, message.Content)
// 	}
// }
