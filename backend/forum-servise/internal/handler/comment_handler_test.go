package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "backend.com/forum/proto"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/forum-servise/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest, opts ...grpc.CallOption) (*pb.ValidateTokenResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.ValidateTokenResponse), args.Error(1)
}

func (m *MockAuthClient) GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.GetUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.GetUserResponse), args.Error(1)
}

func (m *MockAuthClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.LoginResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.LoginResponse), args.Error(1)
}

func (m *MockAuthClient) Register(ctx context.Context, in *pb.RegisterRequest, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.RegisterResponse), args.Error(1)
}

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetCommentsByPostID(ctx context.Context, postID int64) ([]entity.Comment, error) { // Изменен тип возвращаемого значения
	args := m.Called(ctx, postID)
	return args.Get(0).([]entity.Comment), args.Error(1)
}

func TestCreateComment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authClient := new(MockAuthClient)
	commentRepo := new(MockCommentRepository)

	uc := &usecase.CommentUseCase{
		AuthClient:  authClient,
		CommentRepo: commentRepo,
	}

	handler := NewCommentHandler(uc)
	router := gin.Default()
	router.POST("/posts/:id/comments", handler.CreateComment)

	authClient.On("ValidateToken", mock.Anything, &pb.ValidateTokenRequest{Token: "valid-token"}, mock.Anything).
		Return(&pb.ValidateTokenResponse{Valid: true, UserId: 42}, nil)

	authClient.On("GetUser", mock.Anything, &pb.GetUserRequest{Id: 42}, mock.Anything).
		Return(&pb.GetUserResponse{User: &pb.User{Username: "alice"}}, nil)

	expectedComment := &entity.Comment{
		Content:    "test comment",
		AuthorID:   42,
		PostID:     1,
		AuthorName: "alice",
	}
	commentRepo.On("CreateComment", mock.Anything, expectedComment).Return(nil)

	body := `{"content":"test comment"}`
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	authClient.AssertExpectations(t)
	commentRepo.AssertExpectations(t)
}
