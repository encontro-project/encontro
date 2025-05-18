package repository

import (
	"context"

	"encontro/internal/domain/entity"
)

// RoomRepository определяет интерфейс для работы с хранилищем комнат
type RoomRepository interface {
	// CreateRoom создает новую комнату
	CreateRoom(ctx context.Context, room *entity.Room) (*entity.Room, error)
	// GetRoom возвращает комнату по ID
	GetRoom(ctx context.Context, id string) (*entity.Room, error)
	// ListRooms возвращает список всех комнат
	ListRooms(ctx context.Context) ([]*entity.Room, error)
	// UpdateRoom обновляет существующую комнату
	UpdateRoom(ctx context.Context, room *entity.Room) error
	// DeleteRoom удаляет комнату по ID
	DeleteRoom(ctx context.Context, id string) error
	// GetRooms получает пагинированный список комнат
	GetRooms(ctx context.Context, params entity.PaginationParams) ([]*entity.Room, int64, error)
}
