package entity

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID     string
	Conn   *websocket.Conn
	RoomID string
	Send   chan *Envelope
	mu     sync.Mutex
}

// NewClient creates a new client instance
func NewClient(id, roomID string, conn *websocket.Conn) *Client {
	return &Client{
		ID:     id,
		Conn:   conn,
		RoomID: roomID,
		Send:   make(chan *Envelope, 256),
	}
}
