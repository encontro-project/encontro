package usecase

import (
	"context"
	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"encontro/internal/domain/service"
	"time"
)

// MessageUseCase реализует бизнес-логику для работы с сообщениями
type MessageUseCase struct {
	messageRepo repository.MessageRepository
	uuidGen     service.UUIDGenerator
}

// NewMessageUseCase создает новый экземпляр MessageUseCase
func NewMessageUseCase(messageRepo repository.MessageRepository, uuidGen service.UUIDGenerator) *MessageUseCase {
	return &MessageUseCase{
		messageRepo: messageRepo,
		uuidGen:     uuidGen,
	}
}

// CreateMessage создает новое сообщение
func (uc *MessageUseCase) CreateMessage(ctx context.Context, roomID, content, senderID string) (*entity.Message, error) {
	msg := &entity.Message{
		ID:        uc.uuidGen.Generate(),
		RoomID:    roomID,
		Content:   content,
		SenderID:  senderID,
		CreatedAt: time.Now(),
	}

	if err := uc.messageRepo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

// GetMessageByID возвращает сообщение по ID
func (uc *MessageUseCase) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	return uc.messageRepo.GetMessageByID(ctx, id)
}

// GetMessagesByRoomID возвращает пагинированный список сообщений для комнаты
func (uc *MessageUseCase) GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Message], error) {
	messages, total, err := uc.messageRepo.GetMessagesByRoomID(ctx, roomID, params)
	if err != nil {
		return entity.PaginatedResponse[*entity.Message]{}, err
	}

	return entity.NewPaginatedResponse(messages, total, params), nil
}

// DeleteMessage удаляет сообщение по ID
func (uc *MessageUseCase) DeleteMessage(ctx context.Context, id string) error {
	return uc.messageRepo.DeleteMessage(ctx, id)
}
