// internal/usecase/message_usecase.go
package usecase

import (
	"chat-microservice-go/internal/entity"
	"chat-microservice-go/internal/repository"
)

type MessageUseCase struct {
	repo *repository.MessageRepository
}

func NewMessageUseCase(repo *repository.MessageRepository) *MessageUseCase {
	return &MessageUseCase{repo: repo}
}

func (uc *MessageUseCase) SaveMessage(msg entity.Message) error {
	return uc.repo.SaveMessage(msg)
}

func (uc *MessageUseCase) GetMessages() ([]entity.Message, error) {
	return uc.repo.GetMessages()
}
