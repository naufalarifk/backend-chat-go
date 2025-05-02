package websocket

import (
	"fmt"
	"log"
	"realtime-chat-backend/pkg/models"
	"realtime-chat-backend/pkg/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

var store = NewMessageStore()

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		incoming := Message{
			Type: messageType,
			Body: string(p),
		}
		fmt.Printf("Message Received: %+v\n", incoming)
		c.Pool.Broadcast <- incoming

		storedMessages := models.Message{
			ID:        utils.GenerateRandomID(),
			Sender:    c.ID,
			Content:   incoming.Body,
			RoomId:    "default-room",
			Timestamp: time.Now(),
		}
		store.AddMessage(storedMessages)
		for _, msg := range store.GetLastMessages(10) {
			log.Printf("- [%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.Sender, msg.Content)
		}

	}
}
