package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/internal/usecase"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/pkg/auth"
	"github.com/Mandarinka0707/newRepoGOODarhit/chat-servise/pkg/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketController struct {
	authClient  auth.Client
	chatUseCase usecase.ChatUseCase
	logger      logger.Logger
	hub         *Hub
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan entity.Message
	mu         sync.RWMutex
}

func NewWebSocketController(
	authClient *auth.Client,
	chatUseCase usecase.ChatUseCase,
	l logger.Logger,
) *WebSocketController {
	hub := &Hub{
		clients:    make(map[*websocket.Conn]bool),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan entity.Message, 100),
	}
	go hub.run()
	return &WebSocketController{
		authClient:  *authClient,
		chatUseCase: chatUseCase,
		logger:      l,
		hub:         hub,
	}
}

func (c *WebSocketController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	user, err := c.authClient.ValidateToken(r.Context(), token)
	if err != nil {
		c.logger.Errorf("Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	c.hub.register <- conn
	defer func() {
		c.hub.unregister <- conn
	}()

	// Отправка истории сообщений
	messages, err := c.chatUseCase.GetMessages(context.Background())
	if err == nil {
		type ResponseMessage struct {
			ID        int64     `json:"id"`
			UserID    int64     `json:"user_id"`
			Username  string    `json:"username"`
			Content   string    `json:"content"`
			CreatedAt time.Time `json:"created_at"`
		}

		for _, msg := range messages {
			responseMsg := ResponseMessage{
				ID:        msg.ID,
				UserID:    msg.UserID,
				Username:  msg.Username,
				Content:   msg.Content,
				CreatedAt: msg.CreatedAt,
			}

			if err := conn.WriteJSON(responseMsg); err != nil {
				c.logger.Errorf("Error sending history: %v", err)
			}
		}
	}

	// Обработка входящих сообщений
	for {
		var msg entity.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorf("WebSocket error: %v", err)
			}
			break
		}

		msg.UserID = user.ID
		msg.Username = user.Username
		msg.CreatedAt = time.Now()
		c.logger.Infof("Authenticated user: ID=%d, Username=%s", user.ID, user.Username)
		if _, err := c.chatUseCase.CreateMessage(context.Background(), msg); err != nil {
			c.logger.Errorf("Error saving message: %v", err)
			continue
		}

		c.hub.broadcast <- msg
	}
}

func (c *WebSocketController) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := c.chatUseCase.GetMessages(context.Background())
	if err != nil {
		c.logger.Errorf("Failed to get messages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		c.logger.Errorf("Error encoding response: %v", err)
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if err := client.WriteJSON(msg); err != nil {
					log.Printf("Write error: %v", err)
					h.unregister <- client
				}
			}
			h.mu.RUnlock()
		}
	}
}
