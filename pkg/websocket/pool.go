package websocket

import (
	"encoding/json"
	"fmt"
	"realtime-chat-backend/pkg/models"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan models.Envelope
	Chats      map[string][]models.Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan models.Envelope),
		Chats:      make(map[string][]models.Message),
	}
}

func (pool *Pool) Start() {

	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(models.SysMsg{Message: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(models.SysMsg{Message: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			if message.Type == 1 {
				var chat models.Message
				if err := json.Unmarshal(message.Body, &chat); err == nil {
					pool.Chats[chat.RoomId] = append(pool.Chats[chat.RoomId], chat)
				}
				for client := range pool.Clients {
					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println("Error sending to client:", err)
						client.Conn.Close()
						delete(pool.Clients, client)
					}
				}
			}
		}
	}
}
