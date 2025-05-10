// internal/handler/message_handler.go
package handler

import (
	"chat-microservice-go/internal/entity"
	"chat-microservice-go/internal/usecase"
	myWeb "chat-microservice-go/pkg/websocket"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	Uc usecase.MessageUseCase
}

func NewMessageHandler(uc usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{Uc: uc}
}

func (h *MessageHandler) HandleConnections(c *gin.Context) {
	ws, err := myWeb.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	myWeb.Clients[ws] = true

	for {
		var msg entity.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(myWeb.Clients, ws)
			break
		}
		h.Uc.SaveMessage(msg)
		myWeb.Broadcast <- msg
	}
}

func (h *MessageHandler) HandleMessages() {
	for {
		msg := <-myWeb.Broadcast
		for client := range myWeb.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(myWeb.Clients, client)
			}
		}
	}
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	messages, err := h.Uc.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}
