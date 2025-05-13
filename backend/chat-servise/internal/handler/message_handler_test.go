package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat/internal/entity"
	myWeb "github.com/Mandarinka0707/newRepoGOODarhit/chat/pkg/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageUseCase struct {
	mock.Mock
}

func (m *MockMessageUseCase) SaveMessage(msg entity.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageUseCase) GetMessages() ([]entity.Message, error) {
	args := m.Called()
	return args.Get(0).([]entity.Message), args.Error(1)
}

func TestMessageHandler_GetMessages(t *testing.T) {

	uc := new(MockMessageUseCase)

	uc.On("GetMessages").Return([]entity.Message{
		{ID: 1, Username: "testuser", Message: "Hello, World!"},
	}, nil)

	handler := NewMessageHandler(uc)

	router := gin.Default()
	router.GET("/messages", handler.GetMessages)

	req, _ := http.NewRequest("GET", "/messages", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []entity.Message
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, len(resp))
	assert.Equal(t, "Hello, World!", resp[0].Message)
	assert.Equal(t, "testuser", resp[0].Username)
	uc.AssertExpectations(t)
}

func TestMessageHandler_HandleConnections(t *testing.T) {

	uc := new(MockMessageUseCase)

	uc.On("SaveMessage", mock.Anything).Return(nil)

	handler := NewMessageHandler(uc)

	router := gin.Default()
	router.GET("/ws", handler.HandleConnections)

	server := httptest.NewServer(router)
	defer server.Close()

	url := "ws" + server.URL[4:] + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	msg := entity.Message{Username: "testuser", Message: "Hello, World!"}
	err = ws.WriteJSON(msg)
	assert.NoError(t, err)

	uc.AssertExpectations(t)
}

func TestMessageHandler_HandleMessages(t *testing.T) {
	// Создаем мок usecase
	uc := new(MockMessageUseCase)

	// Создаем MessageHandler
	handler := NewMessageHandler(uc)

	// Создаем канал для тестирования
	broadcast := make(chan entity.Message)
	myWeb.Broadcast = broadcast

	// Запускаем горутину для обработки сообщений
	go handler.HandleMessages()

	// Отправляем сообщение в канал
	msg := entity.Message{Username: "testuser", Message: "Hello, World!"}
	broadcast <- msg

	// Проверяем, что сообщение было отправлено всем клиентам
	// (здесь можно добавить проверку для клиентов, если они есть)
}

func TestMessageHandler_GetMessages_Error(t *testing.T) {
	uc := new(MockMessageUseCase)
	uc.On("GetMessages").Return(nil, errors.New("database error"))

	handler := NewMessageHandler(uc)
	router := gin.Default()
	router.GET("/messages", handler.GetMessages)

	req, _ := http.NewRequest("GET", "/messages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}
