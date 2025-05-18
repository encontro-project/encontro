package repository

import (
	"encontro/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

// InMemoryConfig содержит настройки для in-memory режима
type InMemoryConfig struct {
	// InitTestData определяет, нужно ли инициализировать тестовыми данными
	InitTestData bool
}

// NewInMemoryConfig создает новую конфигурацию для in-memory режима
func NewInMemoryConfig(initTestData bool) *InMemoryConfig {
	return &InMemoryConfig{
		InitTestData: initTestData,
	}
}

// generateTestData создает тестовые данные для in-memory режима
func generateTestData() ([]*entity.Room, []*entity.Message) {
	// Создаем тестовые комнаты
	rooms := []*entity.Room{
		{
			ID:        uuid.New().String(),
			Name:      "Общий чат",
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        uuid.New().String(),
			Name:      "Разработка",
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
	}

	// Создаем тестовые сообщения
	messages := []*entity.Message{
		{
			ID:        uuid.New().String(),
			RoomID:    rooms[0].ID,
			Content:   "Привет всем!",
			SenderID:  "user1",
			CreatedAt: time.Now().Add(-23 * time.Hour),
			UpdatedAt: time.Now().Add(-23 * time.Hour),
		},
		{
			ID:        uuid.New().String(),
			RoomID:    rooms[0].ID,
			Content:   "Как дела?",
			SenderID:  "user2",
			CreatedAt: time.Now().Add(-22 * time.Hour),
			UpdatedAt: time.Now().Add(-22 * time.Hour),
		},
		{
			ID:        uuid.New().String(),
			RoomID:    rooms[1].ID,
			Content:   "Нужно исправить баг",
			SenderID:  "user1",
			CreatedAt: time.Now().Add(-11 * time.Hour),
			UpdatedAt: time.Now().Add(-11 * time.Hour),
		},
	}

	return rooms, messages
}
