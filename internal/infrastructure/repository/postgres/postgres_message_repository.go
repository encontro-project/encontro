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

// Create сохраняет новое сообщение
func (r *PostgresMessageRepository) Create(ctx context.Context, msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, room_id, content, sender_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, room_id, content, sender_id, created_at, updated_at`

	msg.UpdatedAt = msg.CreatedAt // Устанавливаем updated_at равным created_at при создании

	err := r.db.WithTransaction(ctx, func(ctx context.Context) error {
		row := r.db.GetPool().QueryRow(ctx, query,
			msg.ID,
			msg.RoomID,
			msg.Content,
			msg.SenderID,
			msg.CreatedAt,
		)

		var returnedMsg entity.Message
		if err := row.Scan(
			&returnedMsg.ID,
			&returnedMsg.RoomID,
			&returnedMsg.Content,
			&returnedMsg.SenderID,
			&returnedMsg.CreatedAt,
			&returnedMsg.UpdatedAt,
		); err != nil {
			return fmt.Errorf("ошибка сканирования созданного сообщения: %w", err)
		}

		// Обновляем поля сообщения из базы данных
		msg.RoomID = returnedMsg.RoomID
		msg.SenderID = returnedMsg.SenderID
		msg.UpdatedAt = returnedMsg.UpdatedAt

		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка создания сообщения: %w", err)
	}

	return nil
}

// GetByID возвращает сообщение по ID
func (r *PostgresMessageRepository) GetByID(ctx context.Context, id string) (*entity.Message, error) {
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

// ListByRoomID возвращает список сообщений в комнате с пагинацией
func (r *PostgresMessageRepository) ListByRoomID(ctx context.Context, roomID string, page, pageSize int) ([]*entity.Message, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Получаем пагинированный список сообщений
	query := `
		SELECT id, room_id, content, sender_id, created_at, updated_at
		FROM messages
		WHERE room_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	offset := (page - 1) * pageSize
	rows, err := r.db.GetPool().Query(ctx, query, roomID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка сообщений: %w", err)
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
			return nil, fmt.Errorf("ошибка сканирования сообщения: %w", err)
		}
		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по сообщениям: %w", err)
	}

	return messages, nil
}

// Delete удаляет сообщение по ID
func (r *PostgresMessageRepository) Delete(ctx context.Context, id string) error {
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

// DeleteByRoomID удаляет все сообщения в комнате
func (r *PostgresMessageRepository) DeleteByRoomID(ctx context.Context, roomID string) error {
	query := `DELETE FROM messages WHERE room_id = $1`

	_, err := r.db.GetPool().Exec(ctx, query, roomID)
	if err != nil {
		return fmt.Errorf("ошибка удаления сообщений комнаты: %w", err)
	}

	return nil
}

// Update обновляет сообщение
func (r *PostgresMessageRepository) Update(ctx context.Context, message *entity.Message) error {
	query := `
		UPDATE messages 
		SET content = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, room_id, content, sender_id, created_at, updated_at`

	err := r.db.GetPool().QueryRow(ctx, query,
		message.Content,
		message.UpdatedAt,
		message.ID,
	).Scan(
		&message.ID,
		&message.RoomID,
		&message.Content,
		&message.SenderID,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("сообщение с ID %s не найдено", message.ID)
	}
	if err != nil {
		return fmt.Errorf("ошибка обновления сообщения: %w", err)
	}

	return nil
}

// CountByRoomID возвращает общее количество сообщений в комнате
func (r *PostgresMessageRepository) CountByRoomID(ctx context.Context, roomID string) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE room_id = $1`
	var count int64
	err := r.db.GetPool().QueryRow(ctx, query, roomID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчета сообщений: %w", err)
	}
	return count, nil
}
