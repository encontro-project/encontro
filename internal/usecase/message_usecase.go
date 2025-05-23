package usecase

import (
	"context"
	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"encontro/internal/domain/service"
	"errors"
	"time"
)

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrMessageNotFound = errors.New("message not found")
	ErrRoomNotFound    = errors.New("room not found")
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
	if roomID == "" || content == "" || senderID == "" {
		return nil, ErrInvalidInput
	}

	msg := &entity.Message{
		ID:        uc.uuidGen.Generate(),
		RoomID:    roomID,
		Content:   content,
		SenderID:  senderID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// Получаем сообщение из БД по ID, чтобы вернуть актуальные значения всех полей
	return uc.messageRepo.GetByID(ctx, msg.ID)
}

// GetMessageByID возвращает сообщение по ID
func (uc *MessageUseCase) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	msg, err := uc.messageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, ErrMessageNotFound
	}
	return msg, nil
}

// GetMessagesByRoomID возвращает пагинированный список сообщений для комнаты
func (uc *MessageUseCase) GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Message], error) {
	if roomID == "" {
		return entity.PaginatedResponse[*entity.Message]{}, ErrInvalidInput
	}

	messages, err := uc.messageRepo.ListByRoomID(ctx, roomID, params.Page, params.PageSize)
	if err != nil {
		return entity.PaginatedResponse[*entity.Message]{}, err
	}

	total, err := uc.messageRepo.CountByRoomID(ctx, roomID)
	if err != nil {
		return entity.PaginatedResponse[*entity.Message]{}, err
	}

	return entity.NewPaginatedResponse(messages, total, params), nil
}

// UpdateMessage обновляет существующее сообщение
func (uc *MessageUseCase) UpdateMessage(ctx context.Context, msg *entity.Message) error {
	if msg == nil || msg.ID == "" || msg.Content == "" {
		return ErrInvalidInput
	}

	// Проверяем существование сообщения
	existingMsg, err := uc.messageRepo.GetByID(ctx, msg.ID)
	if err != nil {
		return err
	}
	if existingMsg == nil {
		return ErrMessageNotFound
	}

	// Обновляем только контент и время обновления
	existingMsg.Content = msg.Content
	existingMsg.UpdatedAt = time.Now()

	return uc.messageRepo.Update(ctx, existingMsg)
}

// DeleteMessage удаляет сообщение по ID
func (uc *MessageUseCase) DeleteMessage(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	// Проверяем существование сообщения
	msg, err := uc.messageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if msg == nil {
		return ErrMessageNotFound
	}

	return uc.messageRepo.Delete(ctx, id)
}
