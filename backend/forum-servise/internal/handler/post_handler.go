package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend.com/forum/forum-servise/internal/repository"
	"backend.com/forum/forum-servise/internal/usecase"
	"backend.com/forum/forum-servise/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostHandler struct {
	uc     *usecase.PostUsecase
	logger *logger.Logger
}

func NewPostHandler(uc *usecase.PostUsecase, logger *logger.Logger) *PostHandler {
	return &PostHandler{uc: uc, logger: logger}
}

func (h *PostHandler) CreatePost(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	var request struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	post, err := h.uc.CreatePost(ctx.Request.Context(), token, request.Title, request.Content)
	if err != nil {
		h.logger.Error("Failed to create post", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":      post.ID,
		"message": "Post created successfully",
		"post":    post,
	})
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, authorNames, err := h.uc.GetPosts(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get posts", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get posts",
			"details": err.Error(),
		})
		return
	}

	response := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		response = append(response, gin.H{
			"id":          post.ID,
			"title":       post.Title,
			"content":     post.Content,
			"author_id":   post.AuthorID,
			"author_name": authorNames[int(post.AuthorID)],
			"created_at":  post.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
func (h *PostHandler) DeletePost(ctx *gin.Context) {
	postID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	h.logger.Debug("Attempting to delete post",
		zap.Int64("post_id", postID),
		zap.String("token", token),
	)

	if err := h.uc.DeletePost(ctx.Request.Context(), token, postID); err != nil {
		h.logger.Error("Failed to delete post", err)

		switch {
		case errors.Is(err, repository.ErrPostNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		case errors.Is(err, repository.ErrPermissionDenied):
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
