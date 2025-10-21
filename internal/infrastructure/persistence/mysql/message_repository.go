package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/domain/repository"
)

type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new MySQL message repository
func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	query := `INSERT INTO message (id, sender, content, room_id, timestamp) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.Sender,
		message.Content,
		message.RoomID,
		message.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id string) (*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp FROM message WHERE id = ?`

	var msg entity.Message
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID,
		&msg.Sender,
		&msg.Content,
		&msg.RoomID,
		&msg.Timestamp,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &msg, nil
}

func (r *messageRepository) GetByRoomID(ctx context.Context, roomID string, limit int) ([]*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp 
	          FROM message 
	          WHERE room_id = ? 
	          ORDER BY timestamp DESC 
	          LIMIT ?`

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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) GetRecent(ctx context.Context, limit int) ([]*entity.Message, error) {
	query := `SELECT id, sender, content, room_id, timestamp 
	          FROM message 
	          ORDER BY timestamp DESC 
	          LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM message WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
