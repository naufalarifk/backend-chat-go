package repository

import (
	"context"
	"realtime-chat-backend/internal/domain/entity"
)

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	// Create inserts a new message
	Create(ctx context.Context, message *entity.Message) error

	// GetByID retrieves a message by its ID
	GetByID(ctx context.Context, id string) (*entity.Message, error)

	// GetByRoomID retrieves all messages for a specific room
	GetByRoomID(ctx context.Context, roomID string, limit int) ([]*entity.Message, error)

	// GetRecent retrieves the most recent messages across all rooms
	GetRecent(ctx context.Context, limit int) ([]*entity.Message, error)

	// Delete removes a message by ID
	Delete(ctx context.Context, id string) error
}
