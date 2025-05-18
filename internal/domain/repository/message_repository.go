package repository

import (
	"context"
	"encontro/internal/domain/entity"
)

// MessageRepository определяет интерфейс для работы с сообщениями
type MessageRepository interface {
	// CreateMessage сохраняет новое сообщение и возвращает его (с ID и временными метками).
	CreateMessage(ctx context.Context, msg *entity.Message) error
	// GetMessageByID возвращает сообщение по его ID.
	GetMessageByID(ctx context.Context, id string) (*entity.Message, error)
	// GetMessagesByRoomID возвращает список сообщений по ID комнаты с пагинацией.
	GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) ([]*entity.Message, int64, error)
	// DeleteMessage удаляет сообщение по его ID.
	DeleteMessage(ctx context.Context, id string) error
}
