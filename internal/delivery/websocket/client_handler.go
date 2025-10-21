package websocket

import (
	"context"
	"encoding/json"
	"log"
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/usecase/chat"
	"realtime-chat-backend/internal/usecase/message"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// ClientHandler handles WebSocket client connections
type ClientHandler struct {
	hub            *chat.Hub
	messageUseCase *message.UseCase
}

// NewClientHandler creates a new client handler
func NewClientHandler(hub *chat.Hub, messageUseCase *message.UseCase) *ClientHandler {
	return &ClientHandler{
		hub:            hub,
		messageUseCase: messageUseCase,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (h *ClientHandler) ReadPump(client *entity.Client) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var envelope entity.Envelope
		if err := json.Unmarshal(message, &envelope); err != nil {
			log.Printf("Failed to unmarshal envelope: %v", err)
			continue
		}

		switch envelope.Type {
		case entity.MessageTypeSystem:
			h.handleSystemMessage(client, &envelope)
		case entity.MessageTypeChat:
			h.handleChatMessage(client, &envelope)
		default:
			log.Printf("Unknown message type: %d", envelope.Type)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (h *ClientHandler) WritePump(client *entity.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case envelope, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(envelope); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *ClientHandler) handleSystemMessage(client *entity.Client, envelope *entity.Envelope) {
	var sysMsg entity.SystemMessage
	if err := json.Unmarshal(envelope.Body, &sysMsg); err != nil {
		log.Printf("Failed to unmarshal system message: %v", err)
		return
	}

	log.Printf("System message from %s: %s", client.ID, sysMsg.Message)

	// Broadcast system message to room
	h.hub.Broadcast(client.RoomID, envelope)
}

func (h *ClientHandler) handleChatMessage(client *entity.Client, envelope *entity.Envelope) {
	var msg entity.Message
	if err := json.Unmarshal(envelope.Body, &msg); err != nil {
		log.Printf("Failed to unmarshal chat message: %v", err)
		return
	}

	// Create message through use case with background context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdMsg, err := h.messageUseCase.CreateMessage(
		ctx,
		msg.Sender,
		msg.Content,
		client.RoomID,
	)
	if err != nil {
		log.Printf("Failed to create message: %v", err)
		return
	}

	// Marshal the created message back to envelope
	msgBytes, err := json.Marshal(createdMsg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	broadcastEnvelope := &entity.Envelope{
		Type: entity.MessageTypeChat,
		Body: msgBytes,
	}

	// Broadcast to room
	h.hub.Broadcast(client.RoomID, broadcastEnvelope)
}

// CreateClient creates a new WebSocket client
func CreateClient(conn *websocket.Conn, roomID string) *entity.Client {
	clientID := uuid.New().String()
	return entity.NewClient(clientID, roomID, conn)
}
