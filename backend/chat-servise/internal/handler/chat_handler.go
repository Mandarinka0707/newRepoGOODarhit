package handler

import (
	"backend.com/forum/chat-servise/internal/usecase"
	"backend.com/forum/chat-servise/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	uc     *usecase.ChatUsecase
	logger *logger.Logger
}

func NewChatHandler(uc *usecase.ChatUsecase, logger *logger.Logger) *ChatHandler {
	return &ChatHandler{uc: uc, logger: logger}
}

func (h *ChatHandler) WebSocketHandler(ctx *gin.Context) {
	conn, err := h.uc.Upgrader().Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", err)
		return
	}

	go h.uc.HandleWebSocket(conn)
}
