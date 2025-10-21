package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	CORS     CORSConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional in production)
	_ = godotenv.Load()

	config := &Config{
		Database: DatabaseConfig{
			User:     getEnv("DBUSER", "root"),
			Password: getEnv("DBPASS", ""),
			Host:     getEnv("DBHOST", "127.0.0.1"),
			Port:     getEnv("DBPORT", "3306"),
			Name:     getEnv("DBNAME", "Messages"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{
				getEnv("CORS_ORIGIN", "http://localhost:3000"),
			},
		},
	}

	// Validate required fields
	if config.Database.User == "" {
		return nil, fmt.Errorf("DBUSER is required")
	}

	return config, nil
}

// GetDSN returns the MySQL Data Source Name
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
