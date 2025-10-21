package http

import (
	"net/http"
	"realtime-chat-backend/internal/domain/entity"
	"realtime-chat-backend/internal/usecase/message"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles HTTP requests for messages
type MessageHandler struct {
	messageUseCase *message.UseCase
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageUseCase *message.UseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	Sender  string `json:"sender" binding:"required"`
	Content string `json:"content" binding:"required"`
	RoomID  string `json:"room_id" binding:"required"`
}

// CreateMessage handles POST /messages
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := h.messageUseCase.CreateMessage(c.Request.Context(), req.Sender, req.Content, req.RoomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetMessages handles GET /messages
func (h *MessageHandler) GetMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	roomID := c.Query("room_id")

	var messages []*entity.Message
	if roomID != "" {
		messages, err = h.messageUseCase.GetMessagesByRoom(c.Request.Context(), roomID, limit)
	} else {
		messages, err = h.messageUseCase.GetRecentMessages(c.Request.Context(), limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
