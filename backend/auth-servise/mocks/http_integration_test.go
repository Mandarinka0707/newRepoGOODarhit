package mocks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"backend.com/forum/auth-servise/internal/controller"
	"backend.com/forum/auth-servise/internal/entity"
	"backend.com/forum/auth-servise/internal/usecase"
	"backend.com/forum/auth-servise/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// UserRepositoryMock - мок репозитория пользователей
type UserRepositoryMock struct {
	users      map[int64]*entity.User
	mu         sync.RWMutex
	forceError bool
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{
		users: make(map[int64]*entity.User),
	}
}

func (r *UserRepositoryMock) CreateUser(ctx context.Context, user *entity.User) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.forceError {
		return 0, errors.New("database error")
	}

	user.ID = int64(len(r.users) + 1)
	r.users[user.ID] = user
	return user.ID, nil
}

func (r *UserRepositoryMock) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.forceError {
		return nil, errors.New("user not found")
	}

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (r *UserRepositoryMock) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.forceError {
		return nil, errors.New("user not found")
	}

	if user, ok := r.users[id]; ok {
		return user, nil
	}
	return nil, nil
}

// SessionRepositoryMock - мок репозитория сессий
type SessionRepositoryMock struct {
	sessions   map[string]*entity.Session
	mu         sync.RWMutex
	forceError bool
}

func NewSessionRepositoryMock() *SessionRepositoryMock {
	return &SessionRepositoryMock{
		sessions: make(map[string]*entity.Session),
	}
}

func (r *SessionRepositoryMock) CreateSession(ctx context.Context, session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.forceError {
		return errors.New("session creation failed")
	}

	r.sessions[session.Token] = session
	return nil
}

func (r *SessionRepositoryMock) GetSessionByToken(ctx context.Context, token string) (*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.forceError {
		return nil, errors.New("session not found")
	}

	if session, ok := r.sessions[token]; ok {
		return session, nil
	}
	return nil, nil
}

func setupHTTPAuthController(t *testing.T) (*gin.Engine, *UserRepositoryMock, *SessionRepositoryMock, func()) {
	logger := zaptest.NewLogger(t)

	userRepo := NewUserRepositoryMock()
	sessionRepo := NewSessionRepositoryMock()

	cfg := &auth.Config{
		TokenSecret:     "test-secret",
		TokenExpiration: 24 * time.Hour,
	}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, cfg, logger)
	httpCtrl := controller.NewHTTPAuthController(uc)

	router := gin.Default()
	router.POST("/register", httpCtrl.Register)
	router.POST("/login", httpCtrl.Login)
	router.GET("/users/:id", httpCtrl.GetUser)

	cleanup := func() {
		userRepo.mu.Lock()
		userRepo.users = make(map[int64]*entity.User)
		userRepo.forceError = false
		userRepo.mu.Unlock()

		sessionRepo.mu.Lock()
		sessionRepo.sessions = make(map[string]*entity.Session)
		sessionRepo.forceError = false
		sessionRepo.mu.Unlock()
	}

	return router, userRepo, sessionRepo, cleanup
}

func TestHTTPAuthController_SuccessFlow(t *testing.T) {
	router, _, _, cleanup := setupHTTPAuthController(t)
	defer cleanup()

	username := "testuser_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// 1. Test successful registration
	registerReq := map[string]string{
		"username": username,
		"password": "testpass123",
	}
	registerBody, _ := json.Marshal(registerReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var registerResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	require.NoError(t, err)
	userID := registerResp["user_id"].(float64)
	require.Greater(t, userID, float64(0))

	// 2. Test successful login
	loginReq := map[string]string{
		"username": username,
		"password": "testpass123",
	}
	loginBody, _ := json.Marshal(loginReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var loginResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)
	require.NotEmpty(t, loginResp["token"])
	require.Equal(t, username, loginResp["username"])

	// 3. Test get user with valid token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/"+strconv.Itoa(int(userID)), nil)
	req.Header.Set("Authorization", "Bearer "+loginResp["token"].(string))
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var userResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &userResp)
	require.NoError(t, err)
	require.Equal(t, username, userResp["username"])
}

func TestHTTPAuthController_ErrorCases(t *testing.T) {
	t.Run("Invalid JSON body", func(t *testing.T) {
		router, _, _, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Invalid request", resp["error"])
	})

	t.Run("Missing required fields", func(t *testing.T) {
		router, _, _, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		invalidReq := map[string]string{"username": "test"} // missing password
		body, _ := json.Marshal(invalidReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Invalid request", resp["error"])
	})

	t.Run("Database error on registration", func(t *testing.T) {
		router, userRepo, _, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		userRepo.forceError = true

		validReq := map[string]string{
			"username": "testuser",
			"password": "testpass",
		}
		body, _ := json.Marshal(validReq)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "database error", resp["error"])
	})

	t.Run("Invalid login credentials", func(t *testing.T) {
		router, userRepo, _, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		// First register a user
		userRepo.forceError = false
		registerReq := map[string]string{
			"username": "testuser",
			"password": "rightpass",
		}
		registerBody, _ := json.Marshal(registerReq)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Then try to login with wrong password
		loginReq := map[string]string{
			"username": "testuser",
			"password": "wrongpass",
		}
		loginBody, _ := json.Marshal(loginReq)
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["error"], "invalid credentials")
	})

	t.Run("Get user without token", func(t *testing.T) {
		router, _, _, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "User not found", resp["error"])
	})

	t.Run("Get user with invalid token", func(t *testing.T) {
		router, _, sessionRepo, cleanup := setupHTTPAuthController(t)
		defer cleanup()

		sessionRepo.forceError = true

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "User not found", resp["error"])
	})
}
