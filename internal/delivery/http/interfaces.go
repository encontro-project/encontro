package http

import (
	"context"
	"encontro/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

// RoomUseCaseInterface определяет интерфейс для работы с комнатами
type RoomUseCaseInterface interface {
	CreateRoom(ctx context.Context, name string, roomType string) (*entity.Room, error)
	GetRoom(ctx context.Context, id string) (*entity.Room, error)
	GetRooms(ctx context.Context, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Room], error)
	ListRooms(ctx context.Context) ([]*entity.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	AddClientToRoom(ctx context.Context, roomID string, client *entity.Client) error
	RemoveClientFromRoom(ctx context.Context, roomID, clientID string) error
}

// MessageUseCaseInterface определяет интерфейс для работы с сообщениями
type MessageUseCaseInterface interface {
	CreateMessage(ctx context.Context, roomID, content, senderID string) (*entity.Message, error)
	GetMessageByID(ctx context.Context, id string) (*entity.Message, error)
	GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Message], error)
	UpdateMessage(ctx context.Context, msg *entity.Message) error
	DeleteMessage(ctx context.Context, id string) error
}

// UserSummaryUseCaseInterface определяет интерфейс для работы с пользователями
type UserSummaryUseCaseInterface interface {
	GetUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error)
	CreateUser(ctx context.Context, user *entity.UserInfo) error
	ChangeUser(ctx context.Context, user *entity.UserInfo) error
	RemoveUser(ctx context.Context, userID int64) error
	CreateServer(ctx context.Context, serverID int64, title string) error
	AddUserToServer(ctx context.Context, userID, serverID int64) error
	CreateChat(ctx context.Context, serverID int64, title, chatType string) error
	PostMessage(ctx context.Context, chatID int64, content string) error
	UpdateAvatar(ctx context.Context, userID int64, newURL string) error
	RenameServer(ctx context.Context, serverID int64, newTitle string) error
}

// WebSocketHandlerInterface определяет интерфейс для работы с WebSocket
type WebSocketHandlerInterface interface {
	HandleWebSocket(c *gin.Context)
}
