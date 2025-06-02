package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"realtime-chat-backend/pkg/database"
	"realtime-chat-backend/pkg/routes"
	"realtime-chat-backend/pkg/websocket"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var store *websocket.MessageStore

//// TODO: create new entry in database

func main() {

	//env load

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	//cap connection

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "Messages"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	database.SetDB(db)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	store = websocket.NewMessageStore()

	if len(store.GetLastMessages(10)) == 0 {
		fmt.Println("✅ NewMessageStore created successfully and is empty.")
	} else {
		fmt.Println("❌ Store is not empty, something is wrong.")
	}
	routes.SetupRoutes(db)
}
