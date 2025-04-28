package controller

import (
	"backend.com/forum/chat-servise/internal/usecase"
	"backend.com/forum/chat-servise/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	uc     *usecase.ChatUsecase
	logger *logger.Logger
}

func NewChatController(uc *usecase.ChatUsecase, logger *logger.Logger) *ChatController {
	return &ChatController{uc: uc, logger: logger}
}

func (c *ChatController) WebSocketHandler(ctx *gin.Context) {
	conn, err := c.uc.Upgrader().Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		c.logger.Error("WebSocket upgrade failed", err)
		return
	}

	go c.uc.HandleWebSocket(conn)
}
