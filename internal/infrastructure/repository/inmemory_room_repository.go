package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
)

// InMemoryRoomRepository реализует RoomRepository с использованием in-memory хранилища
type InMemoryRoomRepository struct {
	rooms map[string]*entity.Room
	mu    sync.RWMutex
}

// NewInMemoryRoomRepository создает новый экземпляр InMemoryRoomRepository
func NewInMemoryRoomRepository(cfg *InMemoryConfig) repository.RoomRepository {
	repo := &InMemoryRoomRepository{
		rooms: make(map[string]*entity.Room),
	}

	// Инициализация тестовыми данными, если указано в конфигурации
	if cfg != nil && cfg.InitTestData {
		rooms, _ := generateTestData()
		for _, room := range rooms {
			repo.rooms[room.ID] = room
		}
	}

	return repo
}

// CreateRoom создает новую комнату
func (r *InMemoryRoomRepository) CreateRoom(ctx context.Context, room *entity.Room) (*entity.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[room.ID]; exists {
		return nil, fmt.Errorf("room with ID %s already exists", room.ID)
	}

	r.rooms[room.ID] = room
	return room, nil
}

// GetRoom возвращает комнату по ID
func (r *InMemoryRoomRepository) GetRoom(ctx context.Context, id string) (*entity.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, fmt.Errorf("room with ID %s not found", id)
	}

	return room, nil
}

// ListRooms возвращает список всех комнат
func (r *InMemoryRoomRepository) ListRooms(ctx context.Context) ([]*entity.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rooms := make([]*entity.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// UpdateRoom обновляет существующую комнату
func (r *InMemoryRoomRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[room.ID]; !exists {
		return fmt.Errorf("room with ID %s not found", room.ID)
	}

	r.rooms[room.ID] = room
	return nil
}

// DeleteRoom удаляет комнату по ID
func (r *InMemoryRoomRepository) DeleteRoom(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[id]; !exists {
		return fmt.Errorf("room with ID %s not found", id)
	}

	delete(r.rooms, id)
	return nil
}

// GetRooms возвращает пагинированный список комнат
func (r *InMemoryRoomRepository) GetRooms(ctx context.Context, params entity.PaginationParams) ([]*entity.Room, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Преобразуем map в slice для сортировки
	rooms := make([]*entity.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	// Сортируем по времени создания (от новых к старым)
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].CreatedAt.After(rooms[j].CreatedAt)
	})

	total := int64(len(rooms))

	// Применяем пагинацию
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if start >= len(rooms) {
		return []*entity.Room{}, total, nil
	}
	if end > len(rooms) {
		end = len(rooms)
	}

	return rooms[start:end], total, nil
}
