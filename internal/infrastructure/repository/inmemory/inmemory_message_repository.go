package repository

import (
	"context"
	"encontro/internal/domain/entity"
	"sort"
	"sync"
)

// InMemoryMessageRepository реализует MessageRepository с использованием in-memory хранилища
type InMemoryMessageRepository struct {
	messages map[string]*entity.Message
	mu       sync.RWMutex
}

// NewInMemoryMessageRepository создает новый экземпляр InMemoryMessageRepository
func NewInMemoryMessageRepository(cfg *InMemoryConfig) *InMemoryMessageRepository {
	repo := &InMemoryMessageRepository{
		messages: make(map[string]*entity.Message),
	}

	// Инициализация тестовыми данными, если указано в конфигурации
	if cfg != nil && cfg.InitTestData {
		_, messages := generateTestData()
		for _, msg := range messages {
			repo.messages[msg.ID] = msg
		}
	}

	return repo
}

// Create сохраняет новое сообщение
func (r *InMemoryMessageRepository) Create(ctx context.Context, msg *entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[msg.ID] = msg
	return nil
}

// GetByID возвращает сообщение по ID
func (r *InMemoryMessageRepository) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if msg, exists := r.messages[id]; exists {
		return msg, nil
	}
	return nil, nil
}

// ListByRoomID возвращает список сообщений в комнате с пагинацией
func (r *InMemoryMessageRepository) ListByRoomID(ctx context.Context, roomID string, page, pageSize int) ([]*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Фильтруем сообщения по roomID
	var roomMessages []*entity.Message
	for _, msg := range r.messages {
		if msg.RoomID == roomID {
			roomMessages = append(roomMessages, msg)
		}
	}

	// Сортируем по времени создания (от новых к старым)
	sort.Slice(roomMessages, func(i, j int) bool {
		return roomMessages[i].CreatedAt.After(roomMessages[j].CreatedAt)
	})

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(roomMessages) {
		return []*entity.Message{}, nil
	}

	if end > len(roomMessages) {
		end = len(roomMessages)
	}

	return roomMessages[start:end], nil
}

// Update обновляет сообщение
func (r *InMemoryMessageRepository) Update(ctx context.Context, message *entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingMessage, exists := r.messages[message.ID]
	if !exists {
		return nil
	}

	message.CreatedAt = existingMessage.CreatedAt
	r.messages[message.ID] = message
	return nil
}

// Delete удаляет сообщение
func (r *InMemoryMessageRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.messages, id)
	return nil
}

// DeleteByRoomID удаляет все сообщения в комнате
func (r *InMemoryMessageRepository) DeleteByRoomID(ctx context.Context, roomID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, msg := range r.messages {
		if msg.RoomID == roomID {
			delete(r.messages, id)
		}
	}

	return nil
}

// CountByRoomID возвращает количество сообщений в комнате
func (r *InMemoryMessageRepository) CountByRoomID(ctx context.Context, roomID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, msg := range r.messages {
		if msg.RoomID == roomID {
			count++
		}
	}
	return count, nil
}
