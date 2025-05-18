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

// CreateMessage сохраняет новое сообщение
func (r *InMemoryMessageRepository) CreateMessage(ctx context.Context, msg *entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages[msg.ID] = msg
	return nil
}

// GetMessageByID возвращает сообщение по ID
func (r *InMemoryMessageRepository) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if msg, exists := r.messages[id]; exists {
		return msg, nil
	}
	return nil, nil
}

// GetMessagesByRoomID возвращает пагинированный список сообщений для комнаты
func (r *InMemoryMessageRepository) GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) ([]*entity.Message, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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

	total := int64(len(roomMessages))

	// Применяем пагинацию
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	if start >= len(roomMessages) {
		return []*entity.Message{}, total, nil
	}
	if end > len(roomMessages) {
		end = len(roomMessages)
	}

	return roomMessages[start:end], total, nil
}

// DeleteMessage удаляет сообщение по ID
func (r *InMemoryMessageRepository) DeleteMessage(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.messages, id)
	return nil
}
