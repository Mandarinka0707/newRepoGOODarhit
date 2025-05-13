package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/repository"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPostUsecase struct {
	mock.Mock
}

func (m *mockPostUsecase) CreatePost(ctx context.Context, token, title, content string) (*entity.Post, error) {
	args := m.Called(ctx, token, title, content)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *mockPostUsecase) GetPosts(ctx context.Context) ([]*entity.Post, map[int]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Post), args.Get(1).(map[int]string), args.Error(2)
}

func (m *mockPostUsecase) DeletePost(ctx context.Context, token string, postID int64) error {
	args := m.Called(ctx, token, postID)
	return args.Error(0)
}

func (m *mockPostUsecase) UpdatePost(ctx context.Context, token string, postID int64, title, content string) (*entity.Post, error) {
	args := m.Called(ctx, token, postID, title, content)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func newTestLogger() *logger.Logger {
	return &logger.Logger{SugaredLogger: nil}
}

func TestCreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockPostUsecase)
	handler := NewPostHandler(mockUC, newTestLogger())

	r := gin.Default()
	r.POST("/posts", handler.CreatePost)

	post := &entity.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  42,
		CreatedAt: time.Now(),
	}

	mockUC.On("CreatePost", mock.Anything, "valid-token", "Test Title", "Test Content").
		Return(post, nil)

	body := `{"title":"Test Title", "content":"Test Content"}`
	req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockPostUsecase)
	handler := NewPostHandler(mockUC, newTestLogger())

	r := gin.Default()
	r.GET("/posts", handler.GetPosts)

	mockPosts := []*entity.Post{
		{ID: 1, Title: "Test", Content: "Body", AuthorID: 1, CreatedAt: time.Now()},
	}
	authors := map[int]string{1: "Alice"}

	mockUC.On("GetPosts", mock.Anything).Return(mockPosts, authors, nil)

	req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeletePost_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockPostUsecase)
	handler := NewPostHandler(mockUC, newTestLogger())

	r := gin.Default()
	r.DELETE("/posts/:id", handler.DeletePost)

	mockUC.On("DeletePost", mock.Anything, "valid-token", int64(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeletePost_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockPostUsecase)
	handler := NewPostHandler(mockUC, newTestLogger())

	r := gin.Default()
	r.DELETE("/posts/:id", handler.DeletePost)

	mockUC.On("DeletePost", mock.Anything, "valid-token", int64(2)).Return(repository.ErrPostNotFound)

	req, _ := http.NewRequest(http.MethodDelete, "/posts/2", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUC := new(mockPostUsecase)
	handler := NewPostHandler(mockUC, newTestLogger())

	r := gin.Default()
	r.PUT("/posts/:id", handler.UpdatePost)

	post := &entity.Post{
		ID:        1,
		Title:     "Updated",
		Content:   "Updated content",
		AuthorID:  42,
		CreatedAt: time.Now(),
	}

	mockUC.On("UpdatePost", mock.Anything, "valid-token", int64(1), "Updated", "Updated content").
		Return(post, nil)

	body := `{"title":"Updated", "content":"Updated content"}`
	req, _ := http.NewRequest(http.MethodPut, "/posts/1", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
