package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"backend.com/forum/chat-servise/internal/entity"
	"backend.com/forum/chat-servise/internal/usecase"
	"backend.com/forum/chat-servise/pkg/auth"
	"backend.com/forum/chat-servise/pkg/logger"
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
	authClient auth.Client,
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
		authClient:  authClient,
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

	// Validate token with auth service
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

	// Register client
	c.hub.register <- conn
	defer func() {
		c.hub.unregister <- conn
	}()

	// Send message history
	messages, err := c.chatUseCase.GetMessages(context.Background())
	if err == nil {
		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				c.logger.Errorf("Error sending history: %v", err)
			}
		}
	}

	// Handle incoming messages
	for {
		var msg entity.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorf("WebSocket error: %v", err)
			}
			break
		}

		// Set user info from validated token
		msg.UserID = user.ID
		msg.Username = user.Username
		msg.CreatedAt = time.Now()

		// Save message to database
		if _, err := c.chatUseCase.CreateMessage(context.Background(), msg); err != nil {
			c.logger.Errorf("Error saving message: %v", err)
			continue
		}

		// Broadcast to all clients
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
