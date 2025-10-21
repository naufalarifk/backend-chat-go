package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"realtime-chat-backend/internal/delivery/http"
	wsDelivery "realtime-chat-backend/internal/delivery/websocket"
	"realtime-chat-backend/internal/infrastructure/config"
	"realtime-chat-backend/internal/infrastructure/persistence/memory"
	"realtime-chat-backend/internal/infrastructure/persistence/mysql"
	"realtime-chat-backend/internal/usecase/chat"
	"realtime-chat-backend/internal/usecase/message"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	messageRepo := mysql.NewMessageRepository(db)
	messageCache := memory.NewMessageCache(1000)

	// Initialize use cases
	messageUseCase := message.NewUseCase(messageRepo, messageCache)
	chatHub := chat.NewHub()

	// Start chat hub
	go chatHub.Run()

	// Initialize handlers
	messageHandler := http.NewMessageHandler(messageUseCase)
	wsHandler := wsDelivery.NewHandler(chatHub, messageUseCase)

	// Setup and start server
	router := setupRouter(cfg, messageHandler, wsHandler)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return db, nil
}

func setupRouter(cfg *config.Config, messageHandler *http.MessageHandler, wsHandler *wsDelivery.Handler) *gin.Engine {
	router := gin.Default()

	// Configure trusted proxies
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Printf("Warning: Failed to set trusted proxies: %v", err)
	}

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	router.GET("/ws", wsHandler.ServeWS)
	router.GET("/messages", messageHandler.GetMessages)
	router.POST("/messages", messageHandler.CreateMessage)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
