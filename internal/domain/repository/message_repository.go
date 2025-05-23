package repository

import (
	"context"
	"encontro/internal/domain/entity"
)

// MessageRepository определяет интерфейс для работы с сообщениями
type MessageRepository interface {
	// Create создает новое сообщение
	Create(ctx context.Context, message *entity.Message) error

	// GetByID возвращает сообщение по ID
	GetByID(ctx context.Context, id string) (*entity.Message, error)

	// ListByRoomID возвращает список сообщений в комнате с пагинацией
	ListByRoomID(ctx context.Context, roomID string, page, pageSize int) ([]*entity.Message, error)

	// CountByRoomID возвращает количество сообщений в комнате
	CountByRoomID(ctx context.Context, roomID string) (int64, error)

	// Update обновляет сообщение
	Update(ctx context.Context, message *entity.Message) error

	// Delete удаляет сообщение
	Delete(ctx context.Context, id string) error

	// DeleteByRoomID удаляет все сообщения в комнате
	DeleteByRoomID(ctx context.Context, roomID string) error
}
