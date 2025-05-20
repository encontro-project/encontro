package usecase

import (
	"context"
	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
)

type UserSummaryUseCase struct {
	repo repository.UserSummaryRepository
}

func NewUserSummaryUseCase(repo repository.UserSummaryRepository) *UserSummaryUseCase {
	return &UserSummaryUseCase{repo: repo}
}

// Read
func (uc *UserSummaryUseCase) GetUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error) {
	return uc.repo.FetchUserSummary(ctx, userID)
}

// Write
func (uc *UserSummaryUseCase) CreateUser(ctx context.Context, user *entity.UserInfo) error {
	return uc.repo.AddUser(ctx, user)
}

func (uc *UserSummaryUseCase) ChangeUser(ctx context.Context, user *entity.UserInfo) error {
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserSummaryUseCase) RemoveUser(ctx context.Context, userID int64) error {
	return uc.repo.DeleteUser(ctx, userID)
}

func (uc *UserSummaryUseCase) CreateServer(ctx context.Context, serverID int64, title string) error {
	return uc.repo.AddServer(ctx, serverID, title)
}

func (uc *UserSummaryUseCase) AddUserToServer(ctx context.Context, userID, serverID int64) error {
	return uc.repo.AddUserToServer(ctx, userID, serverID)
}

func (uc *UserSummaryUseCase) CreateChat(ctx context.Context, serverID int64, title, chatType string) error {
	return uc.repo.AddChat(ctx, serverID, title, chatType)
}

func (uc *UserSummaryUseCase) PostMessage(ctx context.Context, chatID int64, content string) error {
	return uc.repo.AddMessage(ctx, chatID, content)
}

func (uc *UserSummaryUseCase) UpdateAvatar(ctx context.Context, userID int64, newURL string) error {
	return uc.repo.UpdateUserAvatar(ctx, userID, newURL)
}

func (uc *UserSummaryUseCase) RenameServer(ctx context.Context, serverID int64, newTitle string) error {
	return uc.repo.UpdateServerTitle(ctx, serverID, newTitle)
}
