package repository

import (
	"context"
	"fmt"
	"time"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"encontro/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

// PostgresRoomRepository реализует интерфейс RoomRepository для PostgreSQL
type PostgresRoomRepository struct {
	db *database.Pool
}

// NewPostgresRoomRepository создает новый экземпляр PostgresRoomRepository
func NewPostgresRoomRepository(db *database.Pool) repository.RoomRepository {
	return &PostgresRoomRepository{
		db: db,
	}
}

// Create создает новую комнату
func (r *PostgresRoomRepository) Create(ctx context.Context, room *entity.Room) error {
	query := `
		INSERT INTO rooms (id, name, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, type, created_at, updated_at`

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		return r.db.GetPool().QueryRow(ctx, query,
			room.ID,
			room.Name,
			room.Type,
			room.CreatedAt,
			room.UpdatedAt,
		).Scan(
			&room.ID,
			&room.Name,
			&room.Type,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
	})

	if err != nil {
		return fmt.Errorf("ошибка создания комнаты: %w", err)
	}

	return nil
}

// GetByID возвращает комнату по ID
func (r *PostgresRoomRepository) GetByID(ctx context.Context, id string) (*entity.Room, error) {
	query := `
		SELECT id, name, type, created_at, updated_at
		FROM rooms
		WHERE id = $1`

	var room entity.Room
	err := r.db.GetPool().QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.Type,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения комнаты: %w", err)
	}

	return &room, nil
}

// GetRooms возвращает пагинированный список комнат
func (r *PostgresRoomRepository) GetRooms(ctx context.Context, params entity.PaginationParams) ([]*entity.Room, int64, error) {
	// Получаем общее количество комнат
	var total int64
	countQuery := `SELECT COUNT(*) FROM rooms`
	err := r.db.GetPool().QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка подсчета комнат: %w", err)
	}

	// Получаем пагинированный список комнат
	query := `
		SELECT id, name, type, created_at, updated_at
		FROM rooms
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	offset := (params.Page - 1) * params.PageSize
	rows, err := r.db.GetPool().Query(ctx, query, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка комнат: %w", err)
	}
	defer rows.Close()

	var rooms []*entity.Room
	for rows.Next() {
		var room entity.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Type,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования комнаты: %w", err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ошибка при итерации по комнатам: %w", err)
	}

	return rooms, total, nil
}

// List возвращает список комнат с пагинацией
func (r *PostgresRoomRepository) List(ctx context.Context, page, pageSize int) ([]*entity.Room, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	query := `
		SELECT id, name, type, created_at, updated_at
		FROM rooms
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	offset := (page - 1) * pageSize
	rows, err := r.db.GetPool().Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка комнат: %w", err)
	}
	defer rows.Close()

	var rooms []*entity.Room
	for rows.Next() {
		var room entity.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Type,
			&room.CreatedAt,
			&room.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ошибка сканирования комнаты: %w", err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по комнатам: %w", err)
	}

	return rooms, nil
}

// Update обновляет существующую комнату
func (r *PostgresRoomRepository) Update(ctx context.Context, room *entity.Room) error {
	query := `
		UPDATE rooms
		SET name = $1, type = $2, updated_at = $3
		WHERE id = $4`

	room.UpdatedAt = time.Now()
	result, err := r.db.GetPool().Exec(ctx, query,
		room.Name,
		room.Type,
		room.UpdatedAt,
		room.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления комнаты: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("комната с ID %s не найдена", room.ID)
	}

	return nil
}

// Delete удаляет комнату по ID
func (r *PostgresRoomRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM rooms WHERE id = $1`

	result, err := r.db.GetPool().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления комнаты: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("комната с ID %s не найдена", id)
	}

	return nil
}

// AddClient добавляет клиента в комнату
func (r *PostgresRoomRepository) AddClient(ctx context.Context, roomID string, client *entity.Client) error {
	query := `
		INSERT INTO room_clients (room_id, client_id, joined_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, client_id) DO NOTHING`

	_, err := r.db.GetPool().Exec(ctx, query, roomID, client.ID, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка добавления клиента в комнату: %w", err)
	}
	return nil
}

// RemoveClient удаляет клиента из комнаты
func (r *PostgresRoomRepository) RemoveClient(ctx context.Context, roomID string, clientID string) error {
	query := `DELETE FROM room_clients WHERE room_id = $1 AND client_id = $2`
	_, err := r.db.GetPool().Exec(ctx, query, roomID, clientID)
	if err != nil {
		return fmt.Errorf("ошибка удаления клиента из комнаты: %w", err)
	}
	return nil
}

// GetClients возвращает список клиентов в комнате
func (r *PostgresRoomRepository) GetClients(ctx context.Context, roomID string) ([]*entity.Client, error) {
	query := `
		SELECT client_id, joined_at
		FROM room_clients
		WHERE room_id = $1`

	rows, err := r.db.GetPool().Query(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка клиентов: %w", err)
	}
	defer rows.Close()

	var clients []*entity.Client
	for rows.Next() {
		var client entity.Client
		var joinedAt time.Time
		if err := rows.Scan(&client.ID, &joinedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования клиента: %w", err)
		}
		clients = append(clients, &client)
	}

	return clients, nil
}
