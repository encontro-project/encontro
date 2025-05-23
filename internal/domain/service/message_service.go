package service

import (
	"context"
	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"errors"
	"time"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

// MessageService представляет сервис для работы с сообщениями
type MessageService struct {
	messageRepo repository.MessageRepository
	roomRepo    repository.RoomRepository
}

// NewMessageService создает новый экземпляр MessageService
func NewMessageService(messageRepo repository.MessageRepository, roomRepo repository.RoomRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
	}
}

// CreateMessage создает новое сообщение
func (s *MessageService) CreateMessage(ctx context.Context, roomID, senderID, content string) (*entity.Message, error) {
	if roomID == "" || senderID == "" || content == "" {
		return nil, ErrInvalidInput
	}

	// Проверяем существование комнаты
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	now := time.Now()
	message := &entity.Message{
		Content:   content,
		RoomID:    roomID,
		SenderID:  senderID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessage возвращает сообщение по ID
func (s *MessageService) GetMessage(ctx context.Context, id string) (*entity.Message, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	message, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

// ListRoomMessages возвращает список сообщений в комнате с пагинацией
func (s *MessageService) ListRoomMessages(ctx context.Context, roomID string, page, pageSize int) ([]*entity.Message, error) {
	if roomID == "" {
		return nil, ErrInvalidInput
	}

	// Проверяем существование комнаты
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.messageRepo.ListByRoomID(ctx, roomID, page, pageSize)
}

// UpdateMessage обновляет сообщение
func (s *MessageService) UpdateMessage(ctx context.Context, message *entity.Message) error {
	if message == nil || message.ID == "" {
		return ErrInvalidInput
	}

	existingMessage, err := s.messageRepo.GetByID(ctx, message.ID)
	if err != nil {
		return err
	}
	if existingMessage == nil {
		return ErrMessageNotFound
	}

	message.UpdatedAt = time.Now()
	return s.messageRepo.Update(ctx, message)
}

// DeleteMessage удаляет сообщение
func (s *MessageService) DeleteMessage(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	message, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if message == nil {
		return ErrMessageNotFound
	}

	return s.messageRepo.Delete(ctx, id)
}

// DeleteRoomMessages удаляет все сообщения в комнате
func (s *MessageService) DeleteRoomMessages(ctx context.Context, roomID string) error {
	if roomID == "" {
		return ErrInvalidInput
	}

	// Проверяем существование комнаты
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}

	return s.messageRepo.DeleteByRoomID(ctx, roomID)
}
