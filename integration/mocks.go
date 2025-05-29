package integration

import (
	"context"
	"encontro/internal/domain/entity"
	"time"

	"github.com/gin-gonic/gin"
)

// MockRoomUseCase реализует интерфейс для тестов
type MockRoomUseCase struct{}

func (m *MockRoomUseCase) CreateRoom(ctx context.Context, name string, roomType string) (*entity.Room, error) {
	return &entity.Room{
		ID:        "test-room-id",
		Name:      name,
		Type:      roomType,
		Clients:   make([]*entity.Client, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockRoomUseCase) GetRoom(ctx context.Context, id string) (*entity.Room, error) {
	return &entity.Room{
		ID:        id,
		Name:      "Test Room",
		Type:      "text",
		Clients:   make([]*entity.Client, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockRoomUseCase) GetRooms(ctx context.Context, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Room], error) {
	rooms := []*entity.Room{
		{
			ID:        "test-room-1",
			Name:      "Test Room 1",
			Type:      "text",
			Clients:   make([]*entity.Client, 0),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return entity.NewPaginatedResponse(rooms, 1, params), nil
}

func (m *MockRoomUseCase) ListRooms(ctx context.Context) ([]*entity.Room, error) {
	return []*entity.Room{
		{
			ID:        "test-room-1",
			Name:      "Test Room 1",
			Type:      "text",
			Clients:   make([]*entity.Client, 0),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (m *MockRoomUseCase) DeleteRoom(ctx context.Context, id string) error {
	return nil
}

func (m *MockRoomUseCase) AddClientToRoom(ctx context.Context, roomID string, client *entity.Client) error {
	return nil
}

func (m *MockRoomUseCase) RemoveClientFromRoom(ctx context.Context, roomID, clientID string) error {
	return nil
}

// MockMessageUseCase реализует интерфейс для тестов
type MockMessageUseCase struct{}

func (m *MockMessageUseCase) CreateMessage(ctx context.Context, roomID, content, senderID string) (*entity.Message, error) {
	return &entity.Message{
		ID:        "test-message-id",
		RoomID:    roomID,
		Content:   content,
		SenderID:  senderID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockMessageUseCase) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	return &entity.Message{
		ID:        id,
		RoomID:    "test-room-id",
		Content:   "Test message",
		SenderID:  "test-sender-id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockMessageUseCase) GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Message], error) {
	messages := []*entity.Message{
		{
			ID:        "test-message-1",
			RoomID:    roomID,
			Content:   "Test message 1",
			SenderID:  "test-sender-id",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return entity.NewPaginatedResponse(messages, 1, params), nil
}

func (m *MockMessageUseCase) UpdateMessage(ctx context.Context, msg *entity.Message) error {
	return nil
}

func (m *MockMessageUseCase) DeleteMessage(ctx context.Context, id string) error {
	return nil
}

// MockUserSummaryUseCase реализует интерфейс для тестов
type MockUserSummaryUseCase struct{}

func (m *MockUserSummaryUseCase) GetUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error) {
	return &entity.UserSummary{
		Servers: []entity.ServerInfo{
			{
				Title: "Test Server",
				VoiceChannels: []entity.VoiceChannel{
					{
						ID:    1,
						Title: "Test Voice Channel",
					},
				},
				TextChats: []entity.TextChat{
					{
						ID:    1,
						Title: "Test Text Chat",
						Messages: []entity.Message{
							{
								ID:        "test-message-1",
								Content:   "Test message",
								RoomID:    "test-room-1",
								SenderID:  "test-sender-1",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
						},
					},
				},
				AssociatedUsers: []entity.UserInfo{
					{
						ID:        userID,
						Username:  "Test User",
						AvatarURL: "https://example.com/avatar.jpg",
					},
				},
			},
		},
	}, nil
}

func (m *MockUserSummaryUseCase) CreateUser(ctx context.Context, user *entity.UserInfo) error {
	return nil
}

func (m *MockUserSummaryUseCase) ChangeUser(ctx context.Context, user *entity.UserInfo) error {
	return nil
}

func (m *MockUserSummaryUseCase) RemoveUser(ctx context.Context, userID int64) error {
	return nil
}

func (m *MockUserSummaryUseCase) CreateServer(ctx context.Context, serverID int64, title string) error {
	return nil
}

func (m *MockUserSummaryUseCase) AddUserToServer(ctx context.Context, userID, serverID int64) error {
	return nil
}

func (m *MockUserSummaryUseCase) CreateChat(ctx context.Context, serverID int64, title, chatType string) error {
	return nil
}

func (m *MockUserSummaryUseCase) PostMessage(ctx context.Context, chatID int64, content string) error {
	return nil
}

func (m *MockUserSummaryUseCase) UpdateAvatar(ctx context.Context, userID int64, newURL string) error {
	return nil
}

func (m *MockUserSummaryUseCase) RenameServer(ctx context.Context, serverID int64, newTitle string) error {
	return nil
}

// MockWebSocketHandler реализует интерфейс для тестов
type MockWebSocketHandler struct{}

func (m *MockWebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Пустая реализация для тестов
}
