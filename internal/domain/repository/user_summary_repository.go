package repository

import (
	"context"
	"encontro/internal/domain/entity"
)

type UserSummaryRepository interface {
	// Read
	FetchUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error)

	// Write
	AddUser(ctx context.Context, user *entity.UserInfo) error
	UpdateUser(ctx context.Context, user *entity.UserInfo) error
	DeleteUser(ctx context.Context, userID int64) error

	AddServer(ctx context.Context, serverID int64, title string) error
	AddUserToServer(ctx context.Context, userID, serverID int64) error

	AddChat(ctx context.Context, serverID int64, title, chatType string) error
	AddMessage(ctx context.Context, chatID int64, content string) error

	UpdateUserAvatar(ctx context.Context, userID int64, newURL string) error
	UpdateServerTitle(ctx context.Context, serverID int64, newTitle string) error
}
