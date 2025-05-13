// controller/auth_http.go
package controller

import (
	"net/http"
	"strconv"

	"github.com/Mandarinka0707/newRepoGOODarhit/internal/usecase"
	"github.com/gin-gonic/gin"
)

type HTTPAuthController struct {
	uc usecase.AuthUsecaseInterface
}

func NewHTTPAuthController(uc usecase.AuthUsecaseInterface) *HTTPAuthController {
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
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body HTTPRegisterRequest true "Данные для регистрации"
// @Success 200 {object} map[string]interface{} "user_id"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/auth/register [post]
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

	c.JSON(http.StatusOK, gin.H{
		"user_id": ucResp.UserID,
	})
}

// Login выполняет аутентификацию пользователя
// @Summary Аутентификация пользователя
// @Description Вход в систему с логином и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param request body HTTPLoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "token"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/auth/login [post]
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

	c.JSON(http.StatusOK, gin.H{
		"token":    ucResp.Token,
		"username": ucResp.Username,
	})
}

// GetUser получает информацию о пользователе
// @Summary Получить данные пользователя
// @Description Возвращает информацию о пользователе по ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} map[string]interface{} "Данные пользователя"
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/auth/user/{id} [get]
func (ctrl *HTTPAuthController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := ctrl.uc.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
