package handler

import (
	"context"
	"net/http"

	"backend.com/forum/auth-servise/internal/usecase"
	pb "backend.com/forum/proto" // Убрали неиспользуемый импорт logger

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	logger      *zap.Logger
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      logger,
	}
}

func (h *AuthHandler) handleRegister(c *gin.Context) {
	var pbReq pb.RegisterRequest
	if err := c.ShouldBindJSON(&pbReq); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем protobuf запрос в usecase request
	ucReq := &usecase.RegisterRequest{
		Username: pbReq.Username,
		Password: pbReq.Password,
	}

	h.logger.Info("Register request received", zap.String("username", ucReq.Username))

	resp, err := h.authUsecase.Register(context.Background(), ucReq) // Используем правильный тип
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Исправляем имя поля UserID вместо UserId
	h.logger.Info("User registered successfully", zap.Int64("user_id", resp.UserID))
	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserID})
}

func (h *AuthHandler) handleLogin(c *gin.Context) {
	var pbReq pb.LoginRequest
	if err := c.ShouldBindJSON(&pbReq); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем protobuf запрос в usecase request
	ucReq := &usecase.LoginRequest{
		Username: pbReq.Username,
		Password: pbReq.Password,
	}

	resp, err := h.authUsecase.Login(context.Background(), ucReq)
	if err != nil {
		h.logger.Error("Failed to login user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.handleRegister)
		authGroup.POST("/login", h.handleLogin)
	}
}
