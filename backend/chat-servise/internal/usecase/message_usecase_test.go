// chat-microservice-go/internal/usecase/message_usecase_test.go
package usecase

import (
	"errors"
	"testing"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) SaveMessage(msg entity.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessages() ([]entity.Message, error) {
	args := m.Called()
	return args.Get(0).([]entity.Message), args.Error(1)
}
func TestMessageUseCase_SaveMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	uc := NewMessageUseCase(mockRepo)

	mockRepo.On("SaveMessage", entity.Message{Username: "test", Message: "hello"}).Return(nil)
	err := uc.SaveMessage(entity.Message{Username: "test", Message: "hello"})
	assert.NoError(t, err)

	mockRepo.On("SaveMessage", entity.Message{Username: "error", Message: "fail"}).Return(errors.New("db error"))
	err = uc.SaveMessage(entity.Message{Username: "error", Message: "fail"})
	assert.Error(t, err)
}

func TestMessageUseCase_GetMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	uc := NewMessageUseCase(mockRepo)

	expected := []entity.Message{{ID: 1, Username: "user", Message: "test"}}
	mockRepo.On("GetMessages").Return(expected, nil)
	result, err := uc.GetMessages()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.On("GetMessages").Return([]entity.Message{}, nil)
	result, err = uc.GetMessages()
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// func TestMessageUseCase_SaveMessage(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		input       entity.Message
// 		mockSetup   func(*MockMessageRepository)
// 		expectedErr error
// 	}{
// 		{
// 			name: "successful save",
// 			input: entity.Message{
// 				Username: "user1",
// 				Message:  "Hello",
// 			},
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("SaveMessage", entity.Message{
// 					Username: "user1",
// 					Message:  "Hello",
// 				}).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "save error",
// 			input: entity.Message{
// 				Username: "user1",
// 				Message:  "Hello",
// 			},
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("SaveMessage", entity.Message{
// 					Username: "user1",
// 					Message:  "Hello",
// 				}).Return(errors.New("save error"))
// 			},
// 			expectedErr: errors.New("save error"),
// 		},
// 		{
// 			name: "empty message",
// 			input: entity.Message{
// 				Username: "user1",
// 				Message:  "",
// 			},
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("SaveMessage", entity.Message{
// 					Username: "user1",
// 					Message:  "",
// 				}).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "empty username",
// 			input: entity.Message{
// 				Username: "",
// 				Message:  "test",
// 			},
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("SaveMessage", entity.Message{
// 					Username: "",
// 					Message:  "test",
// 				}).Return(nil)
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(MockMessageRepository)
// 			tt.mockSetup(mockRepo)

// 			uc := NewMessageUseCase(mockRepo)
// 			err := uc.SaveMessage(tt.input)

// 			if tt.expectedErr != nil {
// 				assert.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			mockRepo.AssertExpectations(t)
// 		})
// 	}
// }

// func TestMessageUseCase_GetMessages(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		mockSetup    func(*MockMessageRepository)
// 		expectedMsgs []entity.Message
// 		expectedErr  error
// 	}{
// 		{
// 			name: "successful get",
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("GetMessages").Return([]entity.Message{
// 					{Username: "user1", Message: "Hello"},
// 					{Username: "user2", Message: "Hi"},
// 				}, nil)
// 			},
// 			expectedMsgs: []entity.Message{
// 				{Username: "user1", Message: "Hello"},
// 				{Username: "user2", Message: "Hi"},
// 			},
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "get error",
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("GetMessages").Return([]entity.Message{}, errors.New("get error"))
// 			},
// 			expectedMsgs: []entity.Message{},
// 			expectedErr:  errors.New("get error"),
// 		},
// 		{
// 			name: "empty messages",
// 			mockSetup: func(m *MockMessageRepository) {
// 				m.On("GetMessages").Return([]entity.Message{}, nil)
// 			},
// 			expectedMsgs: []entity.Message{},
// 			expectedErr:  nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := new(MockMessageRepository)
// 			tt.mockSetup(mockRepo)

// 			uc := NewMessageUseCase(mockRepo)
// 			messages, err := uc.GetMessages()

// 			assert.Equal(t, tt.expectedMsgs, messages)
// 			if tt.expectedErr != nil {
// 				assert.EqualError(t, err, tt.expectedErr.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			mockRepo.AssertExpectations(t)
// 		})
// 	}
// }
