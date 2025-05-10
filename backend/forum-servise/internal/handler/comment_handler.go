// internal/handler/comment_handler.go
package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"backend.com/forum/forum-servise/internal/entity"
	"backend.com/forum/forum-servise/internal/usecase"
	pb "backend.com/forum/proto"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUC *usecase.CommentUseCase
}

func NewCommentHandler(commentUC *usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{commentUC: commentUC}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Post ID"
// @Param request body entity.Comment true "Comment data"
// @Success 201 {object} entity.Comment
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Аутентификация пользователя
	authResponse, err := h.commentUC.AuthClient.ValidateToken(c.Request.Context(), &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil || authResponse == nil || !authResponse.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Получение post_id из URL
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	// Биндинг данных запроса
	var request struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Получение информации о пользователе
	userResponse, err := h.commentUC.AuthClient.GetUser(c.Request.Context(), &pb.GetUserRequest{
		Id: authResponse.UserId,
	})

	// Создание объекта комментария
	comment := entity.Comment{
		Content:    request.Content,
		AuthorID:   authResponse.UserId,
		AuthorName: "Unknown", // Значение по умолчанию
		PostID:     postID,
	}

	// Обновление имени автора при успешном получении данных
	if err == nil && userResponse != nil && userResponse.User != nil {
		comment.AuthorName = userResponse.User.Username
	}

	// Сохранение комментария
	if err := h.commentUC.CreateComment(c.Request.Context(), &comment); err != nil {
		log.Printf("Error creating comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Возвращаем полные данные комментария
	c.JSON(http.StatusCreated, gin.H{
		"id":          comment.ID,
		"content":     comment.Content,
		"author_id":   comment.AuthorID,
		"post_id":     comment.PostID,
		"author_name": comment.AuthorName,
	})
}

// GetCommentsByPostID godoc
// @Summary Get comments for a post
// @Description Get all comments for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} []entity.Comment
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /api/v1/posts/{id}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	// Получение комментариев
	comments, err := h.commentUC.GetCommentsByPostID(c.Request.Context(), postID)
	if err != nil {
		log.Printf("Error getting comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get comments",
			"details": err.Error(),
		})
		return
	}

	// Правильный формат ответа с массивом комментариев
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
