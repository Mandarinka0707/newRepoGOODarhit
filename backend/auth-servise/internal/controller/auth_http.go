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
		"token": ucResp.Token,
	})
}

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
