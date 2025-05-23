package usecase

import (
	"context"
	"time"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"encontro/internal/domain/service"
)

// RoomUseCase представляет бизнес-логику для работы с комнатами
type RoomUseCase struct {
	repo    repository.RoomRepository
	uuidGen service.UUIDGenerator
}

// NewRoomUseCase создает новый экземпляр RoomUseCase
func NewRoomUseCase(repo repository.RoomRepository, uuidGen service.UUIDGenerator) *RoomUseCase {
	return &RoomUseCase{
		repo:    repo,
		uuidGen: uuidGen,
	}
}

// CreateRoom создает новую комнату
func (uc *RoomUseCase) CreateRoom(ctx context.Context, name string) (*entity.Room, error) {
	now := time.Now()
	room := &entity.Room{
		ID:        uc.uuidGen.Generate(),
		Name:      name,
		Clients:   make([]*entity.Client, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := uc.repo.Create(ctx, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetRoom возвращает комнату по ID
func (uc *RoomUseCase) GetRoom(ctx context.Context, id string) (*entity.Room, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetRooms возвращает пагинированный список комнат
func (uc *RoomUseCase) GetRooms(ctx context.Context, params entity.PaginationParams) (entity.PaginatedResponse[*entity.Room], error) {
	rooms, err := uc.repo.List(ctx, params.Page, params.PageSize)
	if err != nil {
		return entity.PaginatedResponse[*entity.Room]{}, err
	}

	// TODO: получить общее количество комнат
	total := int64(len(rooms))
	return entity.NewPaginatedResponse(rooms, total, params), nil
}

// ListRooms возвращает список всех комнат
func (uc *RoomUseCase) ListRooms(ctx context.Context) ([]*entity.Room, error) {
	return uc.repo.List(ctx, 1, 1000) // TODO: заменить на более элегантное решение
}

// DeleteRoom удаляет комнату по ID
func (uc *RoomUseCase) DeleteRoom(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// AddClientToRoom добавляет клиента в комнату
func (uc *RoomUseCase) AddClientToRoom(ctx context.Context, roomID string, client *entity.Client) error {
	client.ID = uc.uuidGen.Generate()
	return uc.repo.AddClient(ctx, roomID, client)
}

// RemoveClientFromRoom удаляет клиента из комнаты
func (uc *RoomUseCase) RemoveClientFromRoom(ctx context.Context, roomID, clientID string) error {
	return uc.repo.RemoveClient(ctx, roomID, clientID)
}
