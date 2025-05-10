// controller/auth_http.go
package controller

import (
	"net/http"

	"backend.com/forum/auth-servise/internal/usecase"
	"github.com/gin-gonic/gin"
)

type HTTPAuthController struct {
	uc *usecase.AuthUsecase
}

func NewHTTPAuthController(uc *usecase.AuthUsecase) *HTTPAuthController {
	return &HTTPAuthController{uc: uc}
}

type HTTPRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HTTPLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register регистрирует пользователя через HTTP
//
// @Summary Регистрация пользователя через HTTP
// @Description Регистрирует нового пользователя, используя HTTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body HTTPRegisterRequest true "Данные для регистрации"
// @Success 200 {object} map[string]interface{} "Ответ с ID пользователя"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/register [post]
func (ctrl *HTTPAuthController) Register(c *gin.Context) {
	var req HTTPRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ucReq := &usecase.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ucResp, err := ctrl.uc.Register(c.Request.Context(), ucReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Используем gin.H вместо HTTPRegisterResponse
	c.JSON(http.StatusOK, gin.H{
		"user_id": ucResp.UserID,
	})
}

// Login выполняет аутентификацию пользователя через HTTP
//
// @Summary Логин пользователя через HTTP
// @Description Выполняет аутентификацию пользователя, используя HTTP
// @Tags auth
// @Accept json
// @Produce json
// @Param request body HTTPLoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Ответ с токеном и именем пользователя"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/login [post]
func (ctrl *HTTPAuthController) Login(c *gin.Context) {
	var req HTTPLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ucReq := &usecase.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	ucResp, err := ctrl.uc.Login(c.Request.Context(), ucReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Используем gin.H вместо HTTPLoginResponse
	c.JSON(http.StatusOK, gin.H{
		"token":    ucResp.Token,
		"username": ucResp.Username,
	})
}

// GetUser получает информацию о пользователе через HTTP
//
// @Summary Получить информацию о пользователе через HTTP
// @Description Получает информацию о пользователе по ID через HTTP
// @Tags auth
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} map[string]interface{} "Ответ с данными пользователя"
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/user/{id} [get]
func (ctrl *HTTPAuthController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := ctrl.uc.GetUserByID(ctx.Request.Context(), userID)
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
