package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend.com/forum/auth-servise/internal/entity"
	"backend.com/forum/auth-servise/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(ctx context.Context, req *usecase.RegisterRequest) (*usecase.RegisterResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*usecase.RegisterResponse), args.Error(1)
}

func (m *MockAuthUsecase) Login(ctx context.Context, req *usecase.LoginRequest) (*usecase.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*usecase.LoginResponse), args.Error(1)
}

func (m *MockAuthUsecase) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) ValidateToken(ctx context.Context, req *usecase.ValidateTokenRequest) (*usecase.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*usecase.ValidateTokenResponse), args.Error(1)
}

func (m *MockAuthUsecase) GetUser(ctx context.Context, req *usecase.GetUserRequest) (*usecase.GetUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*usecase.GetUserResponse), args.Error(1)
}

func setupRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func TestHandleRegister_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedResp := &usecase.RegisterResponse{UserID: 123}
	mockUsecase.
		On("Register", mock.Anything, &usecase.RegisterRequest{Username: "testuser", Password: "testpass"}).
		Return(expectedResp, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"user_id":123`)
	mockUsecase.AssertExpectations(t)
}

func TestHandleRegister_BindError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestHandleRegister_UsecaseError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockUsecase.
		On("Register", mock.Anything, &usecase.RegisterRequest{Username: "testuser", Password: "testpass"}).
		Return(&usecase.RegisterResponse{}, errors.New("registration failed"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"registration failed"`)
	mockUsecase.AssertExpectations(t)
}

func TestHandleLogin_Success(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedResp := &usecase.LoginResponse{Token: "jwt-token", Username: "testuser"}
	mockUsecase.
		On("Login", mock.Anything, &usecase.LoginRequest{Username: "testuser", Password: "testpass"}).
		Return(expectedResp, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"token":"jwt-token"`)
	mockUsecase.AssertExpectations(t)
}

func TestHandleLogin_BindError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)
	router := setupRouter(h)

	reqBody := map[string]string{
		"username": "baduser",
		"password": "badpass",
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockUsecase.
		On("Login", mock.Anything, &usecase.LoginRequest{Username: "baduser", Password: "badpass"}).
		Return(&usecase.LoginResponse{}, errors.New("invalid credentials"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"invalid credentials"`)
	mockUsecase.AssertExpectations(t)
}

func TestNewAuthHandler(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)

	assert.NotNil(t, h)
	assert.Equal(t, mockUsecase, h.authUsecase)
	assert.Equal(t, logger, h.logger)
}

func TestRegisterRoutes(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	logger := zap.NewNop()
	h := NewAuthHandler(mockUsecase, logger)

	router := gin.New()
	h.RegisterRoutes(router)

	routes := router.Routes()
	assert.Len(t, routes, 2)
	assert.Equal(t, "/auth/register", routes[0].Path)
	assert.Equal(t, "/auth/login", routes[1].Path)
}
