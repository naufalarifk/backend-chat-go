package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	RoomId    string    `json:"roomId"`
	Timestamp time.Time `json:"timestamp"`
}

var messages = []Message{
	{ID: "1", Sender: "Kontos", Content: "Tew", RoomId: "01", Timestamp: time.Now().AddDate(01, 01, 01)},
}

func GetMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, messages)
}
