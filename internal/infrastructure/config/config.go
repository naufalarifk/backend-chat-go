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
	Driver   string //postgres or mysql
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
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
	_ = godotenv.Load()
	config := &Config{
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "postgres"),
			User:     getEnv("DBUSER", "postgres"),
			Password: getEnv("DBPASS", "postgres"),
			Host:     getEnv("DBHOST", "127.0.0.1"),
			Port:     getEnv("DBPORT", "5432"),
			Name:     getEnv("DBNAME", "messages"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
	return config, nil
}

// GetDSN returns the MySQL Data Source Name
func (c *Config) GetDSN() string {
	switch c.Database.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
			c.Database.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	default:
		return ""
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
