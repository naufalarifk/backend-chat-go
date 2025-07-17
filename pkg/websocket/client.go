package websocket

import (
	"encoding/json"
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

var store = NewMessageStore()

func (c *Client) Read() {
	var chatMsg models.Message
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		fmt.Printf("WebSocket frame type: %v\n", messageType)
		if err != nil {
			log.Println(err)
			return
		}

		var env models.Envelope
		log.Println("Raw msg: ", string(p))
		log.Printf("Envelope: %+v", env)
		log.Printf("Body as string: %s", env.Body)

		if err := json.Unmarshal(p, &env); err != nil {
			log.Println("Failed to unmarshal envelope:", err)
			log.Println(p, &env)

			return
		}

		switch env.Type {
		case 0:
			var sysMsg models.SysMsg
			if err := json.Unmarshal(env.Body, &sysMsg); err != nil {
				log.Println("Failed to unmarshal SystemMessage:", err)
				return
			}
			log.Printf(sysMsg.Message)
		case 1:
			var chatMsg models.Message
			if err := json.Unmarshal(env.Body, &chatMsg); err != nil {
				log.Println(env.Body)
				log.Println("Failed to unmarshal ChatMessage:", err)
				return
			}
			log.Printf("Chat from %s: %s", chatMsg.Sender, chatMsg.Content)

		default:
			log.Println("Unknown message type:", env.Type)
		}

		incoming := models.Envelope{
			Type: env.Type,
			Body: env.Body,
		}

		fmt.Printf("Message Received: %+v\n", incoming)
		c.Pool.Broadcast <- incoming

		storedMessages := models.Message{
			ID:        utils.GenerateRandomID(),
			Sender:    chatMsg.Sender,
			Content:   chatMsg.Content,
			RoomId:    chatMsg.RoomId,
			Timestamp: time.Now().String(),
		}

		store.AddMessage(storedMessages)
		for _, msg := range store.GetLastMessages(10) {
			log.Printf("- [%s] %s: %s\n", msg.Timestamp, msg.Sender, msg.Content)
		}

	}
}
