package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/domain/repository"
	"time"

	"github.com/google/uuid"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	query := `INSERT INTO messages (id, sender, content, room_id, timestamp, metadata) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	var id string
	var createdAt time.Time
	msgID := message.ID
	if msgID == "" {
		msgID = uuid.New().String()
		err := r.db.QueryRowContext(ctx, query, msgID, message.Sender, message.Content, message.RoomID, message.Timestamp, "{}").Scan(&id, &createdAt)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

	}
	return nil
}

func (r *messageRepository) GetByRoomID(ctx context.Context, roomID string, limit int) ([]*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp FROM messages WHERE room_id = $1 AND is_deleted = false ORDER BY timestamp DESC LIMIT $2`
	rows, err := r.db.QueryContext(ctx, query, roomID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()
	var messages []*entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.RoomID, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &msg)
	}

	return messages, rows.Err()
}

func (r *messageRepository) Delete(ctx context.Context, messageID string) error {
	query := `UPDATE messages SET is_deleted = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, messageID string) (*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp FROM messages WHERE id = $1 AND is_deleted = false`
	var msg entity.Message
	err := r.db.QueryRowContext(ctx, query, messageID).Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.RoomID, &msg.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}
	return &msg, nil
}

func (r *messageRepository) GetRecent(ctx context.Context, limit int) ([]*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp FROM messages WHERE is_deleted = false ORDER BY timestamp DESC LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
	}
	defer rows.Close()
	var messages []*entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.RoomID, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan recent message: %w", err)
		}
		messages = append(messages, &msg)
	}

	return messages, rows.Err()
}
