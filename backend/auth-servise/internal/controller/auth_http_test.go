// controller/auth_http_test.go
package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mandarinka0707/newRepoGOODarhit/internal/entity"
	"github.com/Mandarinka0707/newRepoGOODarhit/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHTTPAuthController_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockAuthUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful registration",
			requestBody: `{"username": "testuser", "password": "testpass"}`,
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(gomock.Any(), &usecase.RegisterRequest{
					Username: "testuser",
					Password: "testpass",
				}).Return(&usecase.RegisterResponse{UserID: 123}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"user_id":123}`,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"username": "testuser"`, // malformed JSON
			mockSetup:      func(m *MockAuthUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request"}`,
		},
		{
			name:        "usecase error",
			requestBody: `{"username": "testuser", "password": "testpass"}`,
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Register(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"some error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := NewMockAuthUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.Default()
			authController := NewHTTPAuthController(mockUsecase)
			router.POST("/register", authController.Register)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHTTPAuthController_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockAuthUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful login",
			requestBody: `{"username": "testuser", "password": "testpass"}`,
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(gomock.Any(), &usecase.LoginRequest{
					Username: "testuser",
					Password: "testpass",
				}).Return(&usecase.LoginResponse{
					Token:    "testtoken",
					Username: "testuser",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"testtoken","username":"testuser"}`,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"username": "testuser"`, // malformed JSON
			mockSetup:      func(m *MockAuthUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request"}`,
		},
		{
			name:        "usecase error",
			requestBody: `{"username": "testuser", "password": "testpass"}`,
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().Login(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"invalid credentials"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := NewMockAuthUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.Default()
			authController := NewHTTPAuthController(mockUsecase)
			router.POST("/login", authController.Login)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestHTTPAuthController_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockAuthUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful get user",
			userID: "123",
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().GetUserByID(gomock.Any(), int64(123)).
					Return(&entity.User{
						ID:       123,
						Username: "testuser",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":123,"username":"testuser"}`,
		},
		{
			name:   "user not found",
			userID: "456",
			mockSetup: func(m *MockAuthUsecase) {
				m.EXPECT().GetUserByID(gomock.Any(), int64(456)).
					Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"User not found"}`,
		},
		{
			name:   "invalid user ID format",
			userID: "abc",
			mockSetup: func(m *MockAuthUsecase) {
				// No expectation as it should fail before calling usecase
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"User not found"}`,
		},
		{
			name:   "empty user ID",
			userID: "",
			mockSetup: func(m *MockAuthUsecase) {
				// No expectation as it should fail before calling usecase
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"User not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := NewMockAuthUsecase(ctrl)
			tt.mockSetup(mockUsecase)

			router := gin.Default()
			authController := NewHTTPAuthController(mockUsecase)
			router.GET("/user/:id", authController.GetUser)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/user/"+tt.userID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			} else {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
