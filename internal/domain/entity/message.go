package entity

import (
	"encoding/json"
	"time"
)

// Envelope wraps WebSocket messages with type information
type Envelope struct {
	Type int             `json:"type"`
	Body json.RawMessage `json:"body"`
}

// SystemMessage represents system-level messages (joins, disconnects, etc.)
type SystemMessage struct {
	Message string `json:"message"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageType constants
const (
	MessageTypeSystem = 0
	MessageTypeChat   = 1
)
