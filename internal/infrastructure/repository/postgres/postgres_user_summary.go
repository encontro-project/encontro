package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"encontro/internal/domain/entity"
	"encontro/internal/domain/repository"
)

type PostgresUserSummaryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserSummaryRepository(pool *pgxpool.Pool) repository.UserSummaryRepository {
	return &PostgresUserSummaryRepository{pool: pool}
}

// FetchUserSummary loads a UserSummary by userID.
func (r *PostgresUserSummaryRepository) FetchUserSummary(ctx context.Context, userID int64) (*entity.UserSummary, error) {
	// TODO: заполнить SELECT-запросы для получения UserSummary.Servers, VoiceChannels, TextChats, Messages, AssociatedUsers

	// stub чтобы код компилировался
	return &entity.UserSummary{}, nil
}

func (r *PostgresUserSummaryRepository) AddUser(ctx context.Context, user *entity.UserInfo) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, avatar_url, created_at, updated_at)
         VALUES ($1, $2, $3, NOW(), NOW())`,
		user.ID, user.Username, user.AvatarURL,
	)
	return err
}

func (r *PostgresUserSummaryRepository) UpdateUser(ctx context.Context, user *entity.UserInfo) error {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE users
         SET username = $1, avatar_url = $2, updated_at = NOW()
         WHERE id = $3`,
		user.Username, user.AvatarURL, user.ID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %d", user.ID)
	}
	return nil
}

func (r *PostgresUserSummaryRepository) DeleteUser(ctx context.Context, userID int64) error {
	cmdTag, err := r.pool.Exec(ctx,
		`DELETE FROM users WHERE id = $1`, userID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %d", userID)
	}
	return nil
}

func (r *PostgresUserSummaryRepository) AddServer(ctx context.Context, serverID int64, title string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO servers (id, name, created_at, updated_at)
         VALUES ($1, $2, NOW(), NOW())`,
		serverID, title,
	)
	return err
}

func (r *PostgresUserSummaryRepository) AddUserToServer(ctx context.Context, userID, serverID int64) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO server_users (server_id, user_id, created_at)
         VALUES ($1, $2, NOW())`,
		serverID, userID,
	)
	return err
}

func (r *PostgresUserSummaryRepository) AddChat(ctx context.Context, serverID int64, title, chatType string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO chats (server_id, title, type, created_at)
         VALUES ($1, $2, $3, NOW())`,
		serverID, title, chatType,
	)
	return err
}

func (r *PostgresUserSummaryRepository) AddMessage(ctx context.Context, chatID int64, content string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO messages (chat_id, content, created_at)
         VALUES ($1, $2, NOW())`,
		chatID, content,
	)
	return err
}

func (r *PostgresUserSummaryRepository) UpdateUserAvatar(ctx context.Context, userID int64, newURL string) error {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE users
         SET avatar_url = $1, updated_at = NOW()
         WHERE id = $2`,
		newURL, userID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found: %d", userID)
	}
	return nil
}

func (r *PostgresUserSummaryRepository) UpdateServerTitle(ctx context.Context, serverID int64, newTitle string) error {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE servers
         SET name = $1, updated_at = NOW()
         WHERE id = $2`,
		newTitle, serverID,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("server not found: %d", serverID)
	}
	return nil
}
