package websocket

import (
	"log"
	"net/http"
	"realtime-chat-backend/internal/usecase/chat"
	"realtime-chat-backend/internal/usecase/message"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub           *chat.Hub
	clientHandler *ClientHandler
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *chat.Hub, messageUseCase *message.UseCase) *Handler {
	return &Handler{
		hub:           hub,
		clientHandler: NewClientHandler(hub, messageUseCase),
	}
}

// ServeWS handles WebSocket requests from clients
func (h *Handler) ServeWS(c *gin.Context) {
	roomID := c.DefaultQuery("roomId", "default")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := CreateClient(conn, roomID)
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go h.clientHandler.WritePump(client)
	go h.clientHandler.ReadPump(client)
}
