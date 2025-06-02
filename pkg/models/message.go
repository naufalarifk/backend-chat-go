package models

import (
	"database/sql"
	"log"
	"net/http"
	"realtime-chat-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	RoomId    string `json:"roomId"`
	Timestamp string `json:"timestamp"`
}

func GetMessages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var messages []Message

		rows, err := database.DB.Query("SELECT * FROM message")
		if err != nil {
			log.Printf("Query error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query messages"})

		}

		for rows.Next() {
			var msg Message
			if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.RoomId, &msg.Timestamp); err != nil {
				log.Printf("Scan error: %v", err)
				continue
			}
			messages = append(messages, msg)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Row iteration error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read message rows"})
			return
		}

		c.IndentedJSON(http.StatusOK, messages)
	}
}
