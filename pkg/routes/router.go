package routes

import (
	"database/sql"
	"realtime-chat-backend/pkg/models"
	"realtime-chat-backend/pkg/websocket"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) {
	router := gin.Default()
	pool := websocket.NewPool()
	go pool.Start()

	err := router.SetTrustedProxies([]string{"127.0.0.1"})

	if err != nil {
		panic("Error Trusted Proxies")
	}

	//enable CORS

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//SERVE WS

	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(pool, c.Writer, c.Request)
	})

	//SERVE MESSAGES

	router.GET("/messages", models.GetMessages(db))
	router.POST("/messages", models.AddNewMessage(db))

	router.Run(":8080")
}
