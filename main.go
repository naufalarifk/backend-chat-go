package main

import (
	"fmt"
	"realtime-chat-backend/pkg/routes"
	"realtime-chat-backend/pkg/websocket"
)

var store *websocket.MessageStore

func main() {
	store = websocket.NewMessageStore()

	if len(store.GetLastMessages(10)) == 0 {
		fmt.Println("✅ NewMessageStore created successfully and is empty.")
	} else {
		fmt.Println("❌ Store is not empty, something is wrong.")
	}
	routes.SetupRoutes()
}
