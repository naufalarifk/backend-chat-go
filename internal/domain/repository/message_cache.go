package repository

import (
	"realtime-chat-backend/internal/domain/entity"
)

// MessageCache defines the interface for in-memory message caching
type MessageCache interface {
	// Add stores a message in cache
	Add(message *entity.Message)

	// GetRecent retrieves the most recent n messages
	GetRecent(n int) []*entity.Message

	// GetByRoomID retrieves messages for a specific room
	GetByRoomID(roomID string, limit int) []*entity.Message

	// Clear removes all messages from cache
	Clear()

	// Count returns the number of messages in cache
	Count() int
}
