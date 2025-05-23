package repository

import (
	"context"

	"encontro/internal/domain/entity"
)

// RoomRepository определяет интерфейс для работы с комнатами
type RoomRepository interface {
	// Create создает новую комнату
	Create(ctx context.Context, room *entity.Room) error

	// GetByID возвращает комнату по ID
	GetByID(ctx context.Context, id string) (*entity.Room, error)

	// List возвращает список всех комнат с пагинацией
	List(ctx context.Context, page, pageSize int) ([]*entity.Room, error)

	// Update обновляет информацию о комнате
	Update(ctx context.Context, room *entity.Room) error

	// Delete удаляет комнату
	Delete(ctx context.Context, id string) error

	// AddClient добавляет клиента в комнату
	AddClient(ctx context.Context, roomID string, client *entity.Client) error

	// RemoveClient удаляет клиента из комнаты
	RemoveClient(ctx context.Context, roomID, clientID string) error

	// GetClients возвращает список клиентов в комнате
	GetClients(ctx context.Context, roomID string) ([]*entity.Client, error)
}
