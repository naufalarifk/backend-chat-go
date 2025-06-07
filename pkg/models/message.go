package models

import (
	"database/sql"
	"log"
	"net/http"
	"realtime-chat-backend/pkg/database"
	"time"

	"github.com/gin-gonic/gin"
)

//typing and error handling

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

func AddNewMessage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg Message

		if err := c.BindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if msg.Timestamp == "" {
			msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		} else {
			parsed, err := time.Parse(time.RFC3339, msg.Timestamp)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
				return
			}
			msg.Timestamp = parsed.Format("2006-01-02 15:04:05")
		}

		query := "INSERT INTO message (Sender, Content, RoomId, Timestamp) VALUES (?, ?, ?, ?)"

		_, err := database.DB.ExecContext(c, query, msg.Sender, msg.Content, msg.RoomId, msg.Timestamp)

		if err != nil {
			log.Printf("Insert Error: %v:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed To Insert Message"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Message Added Successfully!"})
	}
}
