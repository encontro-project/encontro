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

// CreateRoom создает новую комнату
func (r *PostgresRoomRepository) CreateRoom(ctx context.Context, room *entity.Room) (*entity.Room, error) {
	query := `
		INSERT INTO rooms (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, created_at, updated_at`

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		return r.db.GetPool().QueryRow(ctx, query,
			room.ID,
			room.Name,
			room.CreatedAt,
			room.UpdatedAt,
		).Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка создания комнаты: %w", err)
	}

	return room, nil
}

// GetRoom возвращает комнату по ID
func (r *PostgresRoomRepository) GetRoom(ctx context.Context, id string) (*entity.Room, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM rooms
		WHERE id = $1`

	var room entity.Room
	err := r.db.GetPool().QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
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
		SELECT id, name, created_at, updated_at
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

// ListRooms возвращает список всех комнат
func (r *PostgresRoomRepository) ListRooms(ctx context.Context) ([]*entity.Room, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM rooms
		ORDER BY created_at DESC`

	rows, err := r.db.GetPool().Query(ctx, query)
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

// UpdateRoom обновляет существующую комнату
func (r *PostgresRoomRepository) UpdateRoom(ctx context.Context, room *entity.Room) error {
	query := `
		UPDATE rooms
		SET name = $1, updated_at = $2
		WHERE id = $3`

	room.UpdatedAt = time.Now()
	result, err := r.db.GetPool().Exec(ctx, query,
		room.Name,
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

// DeleteRoom удаляет комнату по ID
func (r *PostgresRoomRepository) DeleteRoom(ctx context.Context, id string) error {
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
