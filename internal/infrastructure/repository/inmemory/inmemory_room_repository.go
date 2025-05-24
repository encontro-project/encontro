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
func NewInMemoryRoomRepository(cfg *InMemoryConfig) (repository.RoomRepository, error) {
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

	return repo, nil
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

// GetByID возвращает комнату по ID
func (r *InMemoryRoomRepository) GetByID(ctx context.Context, id string) (*entity.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, fmt.Errorf("room with ID %s not found", id)
	}

	return room, nil
}

// List возвращает список всех комнат с пагинацией
func (r *InMemoryRoomRepository) List(ctx context.Context, page, pageSize int) ([]*entity.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	rooms := make([]*entity.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	// Сортируем по времени создания (от новых к старым)
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].CreatedAt.After(rooms[j].CreatedAt)
	})

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(rooms) {
		return []*entity.Room{}, nil
	}

	if end > len(rooms) {
		end = len(rooms)
	}

	return rooms[start:end], nil
}

// Update обновляет существующую комнату
func (r *InMemoryRoomRepository) Update(ctx context.Context, room *entity.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[room.ID]; !exists {
		return fmt.Errorf("room with ID %s not found", room.ID)
	}

	r.rooms[room.ID] = room
	return nil
}

// Delete удаляет комнату по ID
func (r *InMemoryRoomRepository) Delete(ctx context.Context, id string) error {
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

// AddClient добавляет клиента в комнату
func (r *InMemoryRoomRepository) AddClient(ctx context.Context, roomID string, client *entity.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return fmt.Errorf("room with ID %s not found", roomID)
	}

	if room.Clients == nil {
		room.Clients = make([]*entity.Client, 0)
	}
	room.Clients = append(room.Clients, client)
	return nil
}

// Create создает новую комнату
func (r *InMemoryRoomRepository) Create(ctx context.Context, room *entity.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rooms[room.ID]; exists {
		return fmt.Errorf("room with ID %s already exists", room.ID)
	}

	r.rooms[room.ID] = room
	return nil
}

// GetClients возвращает список клиентов в комнате
func (r *InMemoryRoomRepository) GetClients(ctx context.Context, roomID string) ([]*entity.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room with ID %s not found", roomID)
	}

	if room.Clients == nil {
		return []*entity.Client{}, nil
	}

	return room.Clients, nil
}

// RemoveClient удаляет клиента из комнаты
func (r *InMemoryRoomRepository) RemoveClient(ctx context.Context, roomID string, clientID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return fmt.Errorf("room with ID %s not found", roomID)
	}

	if room.Clients == nil {
		return nil
	}

	// Ищем и удаляем клиента
	for i, client := range room.Clients {
		if client.ID == clientID {
			room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
			return nil
		}
	}

	return nil
}
