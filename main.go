package main

import (
	"fmt"
	"net/http"

	"realtime-chat-backend/pkg/websocket"
)

var store *websocket.MessageStore

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	store = websocket.NewMessageStore()

	if len(store.GetLastMessages(10)) == 0 {
		fmt.Println("✅ NewMessageStore created successfully and is empty.")
	} else {
		fmt.Println("❌ Store is not empty, something is wrong.")
	}
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
