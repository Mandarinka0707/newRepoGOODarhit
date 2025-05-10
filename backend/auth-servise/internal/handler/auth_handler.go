package handler

import (
	"context"
	"net/http"

	"backend.com/forum/auth-servise/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecaseInterface // Используем интерфейс
	logger      *zap.Logger
}

// AuthUsecaseInterface определяет методы, которые использует handler
type AuthUsecaseInterface interface {
	Register(ctx context.Context, req *usecase.RegisterRequest) (*usecase.RegisterResponse, error)
	Login(ctx context.Context, req *usecase.LoginRequest) (*usecase.LoginResponse, error)
}

func NewAuthHandler(authUsecase usecase.AuthUsecaseInterface, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		logger:      logger,
	}
}

// handleRegister обрабатывает регистрацию пользователя.
//
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} map[string]interface{} "user_id"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) handleRegister(c *gin.Context) {
	var req usecase.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.Register(context.Background(), &req)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserID})
}

// handleLogin обрабатывает авторизацию пользователя.
//
// @Summary Авторизация пользователя
// @Description Логин по имени пользователя и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecase.LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "token"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) handleLogin(c *gin.Context) {
	var req usecase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.Login(context.Background(), &req)
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
