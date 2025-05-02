package routes

import (
	"realtime-chat-backend/pkg/models"
	"realtime-chat-backend/pkg/websocket"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	router := gin.Default()
	pool := websocket.NewPool()
	go pool.Start()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(pool, c.Writer, c.Request)
	})

	router.GET("/messages", models.GetMessages)

	router.Run(":8080")
}
