package chat

import (
	"realtime-chat-backend/internal/domain/entity"
	"sync"
)

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients per room
	rooms map[string]map[*entity.Client]bool

	// Register requests from clients
	register chan *entity.Client

	// Unregister requests from clients
	unregister chan *entity.Client

	// Broadcast messages to all clients in a room
	broadcast chan *BroadcastMessage

	mu sync.RWMutex
}

// BroadcastMessage wraps a message with room information
type BroadcastMessage struct {
	RoomID   string
	Envelope *entity.Envelope
}

// NewHub creates a new chat hub
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*entity.Client]bool),
		register:   make(chan *entity.Client),
		unregister: make(chan *entity.Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToRoom(message)
		}
	}
}

func (h *Hub) registerClient(client *entity.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.RoomID] == nil {
		h.rooms[client.RoomID] = make(map[*entity.Client]bool)
	}
	h.rooms[client.RoomID][client] = true
}

func (h *Hub) unregisterClient(client *entity.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.RoomID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.Send)

			// Clean up empty rooms
			if len(clients) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
	}
}

func (h *Hub) broadcastToRoom(message *BroadcastMessage) {
	h.mu.RLock()
	clients := h.rooms[message.RoomID]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.Send <- message.Envelope:
		default:
			// Client's send channel is full, unregister
			h.unregister <- client
		}
	}
}

// Register registers a client
func (h *Hub) Register(client *entity.Client) {
	h.register <- client
}

// Unregister unregisters a client
func (h *Hub) Unregister(client *entity.Client) {
	h.unregister <- client
}

// Broadcast broadcasts a message to a room
func (h *Hub) Broadcast(roomID string, envelope *entity.Envelope) {
	h.broadcast <- &BroadcastMessage{
		RoomID:   roomID,
		Envelope: envelope,
	}
}

// GetRoomClientCount returns the number of clients in a room
func (h *Hub) GetRoomClientCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}
