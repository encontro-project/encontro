package service

import (
	"context"
	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"errors"
	"time"
)

var (
	ErrRoomNotFound = errors.New("room not found")
	ErrInvalidInput = errors.New("invalid input")
)

// RoomService представляет сервис для работы с комнатами
type RoomService struct {
	roomRepo repository.RoomRepository
}

// NewRoomService создает новый экземпляр RoomService
func NewRoomService(roomRepo repository.RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

// CreateRoom создает новую комнату
func (s *RoomService) CreateRoom(ctx context.Context, name string) (*entity.Room, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}

	room := entity.NewRoom("", name) // ID будет сгенерирован в репозитории
	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

// GetRoom возвращает комнату по ID
func (s *RoomService) GetRoom(ctx context.Context, id string) (*entity.Room, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

// ListRooms возвращает список комнат с пагинацией
func (s *RoomService) ListRooms(ctx context.Context, page, pageSize int) ([]*entity.Room, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return s.roomRepo.List(ctx, page, pageSize)
}

// UpdateRoom обновляет информацию о комнате
func (s *RoomService) UpdateRoom(ctx context.Context, room *entity.Room) error {
	if room == nil || room.ID == "" {
		return ErrInvalidInput
	}

	existingRoom, err := s.roomRepo.GetByID(ctx, room.ID)
	if err != nil {
		return err
	}
	if existingRoom == nil {
		return ErrRoomNotFound
	}

	room.UpdatedAt = time.Now()
	return s.roomRepo.Update(ctx, room)
}

// DeleteRoom удаляет комнату
func (s *RoomService) DeleteRoom(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}

	return s.roomRepo.Delete(ctx, id)
}

// AddClientToRoom добавляет клиента в комнату
func (s *RoomService) AddClientToRoom(ctx context.Context, roomID string, client *entity.Client) error {
	if roomID == "" || client == nil {
		return ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}

	return s.roomRepo.AddClient(ctx, roomID, client)
}

// RemoveClientFromRoom удаляет клиента из комнаты
func (s *RoomService) RemoveClientFromRoom(ctx context.Context, roomID, clientID string) error {
	if roomID == "" || clientID == "" {
		return ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return ErrRoomNotFound
	}

	return s.roomRepo.RemoveClient(ctx, roomID, clientID)
}

// GetRoomClients возвращает список клиентов в комнате
func (s *RoomService) GetRoomClients(ctx context.Context, roomID string) ([]*entity.Client, error) {
	if roomID == "" {
		return nil, ErrInvalidInput
	}

	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, ErrRoomNotFound
	}

	return s.roomRepo.GetClients(ctx, roomID)
}
