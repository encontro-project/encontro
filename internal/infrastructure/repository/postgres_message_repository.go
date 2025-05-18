package repository

import (
	"context"
	"fmt"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
	"encontro/internal/infrastructure/database"

	"github.com/jackc/pgx/v5"
)

// PostgresMessageRepository реализует MessageRepository с использованием PostgreSQL
type PostgresMessageRepository struct {
	db *database.Pool
}

// NewPostgresMessageRepository создает новый экземпляр PostgresMessageRepository
func NewPostgresMessageRepository(db *database.Pool) repository.MessageRepository {
	return &PostgresMessageRepository{
		db: db,
	}
}

// CreateMessage сохраняет новое сообщение
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, room_id, content, sender_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, room_id, content, sender_id, created_at, updated_at`

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		return r.db.GetPool().QueryRow(ctx, query,
			msg.ID,
			msg.RoomID,
			msg.Content,
			msg.SenderID,
			msg.CreatedAt,
			msg.UpdatedAt,
		).Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.Content,
			&msg.SenderID,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
	})

	if err != nil {
		return fmt.Errorf("ошибка создания сообщения: %w", err)
	}

	return nil
}

// GetMessageByID возвращает сообщение по ID
func (r *PostgresMessageRepository) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	query := `
		SELECT id, room_id, content, sender_id, created_at, updated_at
		FROM messages
		WHERE id = $1`

	var msg entity.Message
	err := r.db.GetPool().QueryRow(ctx, query, id).Scan(
		&msg.ID,
		&msg.RoomID,
		&msg.Content,
		&msg.SenderID,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сообщения: %w", err)
	}

	return &msg, nil
}

// GetMessagesByRoomID возвращает пагинированный список сообщений для комнаты
func (r *PostgresMessageRepository) GetMessagesByRoomID(ctx context.Context, roomID string, params entity.PaginationParams) ([]*entity.Message, int64, error) {
	// Получаем общее количество сообщений в комнате
	var total int64
	countQuery := `SELECT COUNT(*) FROM messages WHERE room_id = $1`
	err := r.db.GetPool().QueryRow(ctx, countQuery, roomID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка подсчета сообщений: %w", err)
	}

	// Получаем пагинированный список сообщений
	query := `
		SELECT id, room_id, content, sender_id, created_at, updated_at
		FROM messages
		WHERE room_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	offset := (params.Page - 1) * params.PageSize
	rows, err := r.db.GetPool().Query(ctx, query, roomID, params.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка сообщений: %w", err)
	}
	defer rows.Close()

	var messages []*entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.Content,
			&msg.SenderID,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("ошибка сканирования сообщения: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ошибка при итерации по сообщениям: %w", err)
	}

	return messages, total, nil
}

// DeleteMessage удаляет сообщение по ID
func (r *PostgresMessageRepository) DeleteMessage(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.GetPool().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления сообщения: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("сообщение с ID %s не найдено", id)
	}

	return nil
}
