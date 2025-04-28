// internal/usecase/chat_usecase.go
package usecase

import (
	"context"

	"backend.com/forum/chat-servise/internal/entity"
	"backend.com/forum/chat-servise/internal/repository"
)

type ChatUseCase interface {
	CreateMessage(ctx context.Context, msg entity.Message) (*entity.Message, error)
	GetMessages(ctx context.Context) ([]entity.Message, error)
}

type chatUseCase struct {
	repo repository.MessageRepository
}

func NewChatUsecase(repo repository.MessageRepository) ChatUseCase {
	return &chatUseCase{repo: repo}
}

func (uc *chatUseCase) CreateMessage(ctx context.Context, msg entity.Message) (*entity.Message, error) {
	message := repository.Message{
		UserID:    msg.UserID,
		Username:  msg.Username,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}

	created, err := uc.repo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	return &entity.Message{
		ID:        created.ID,
		UserID:    created.UserID,
		Username:  created.Username,
		Content:   created.Content,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (uc *chatUseCase) GetMessages(ctx context.Context) ([]entity.Message, error) {
	messages, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Message, len(messages))
	for i, m := range messages {
		result[i] = entity.Message{
			ID:        m.ID,
			UserID:    m.UserID,
			Username:  m.Username,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		}
	}

	return result, nil
}
