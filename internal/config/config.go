package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config represents application configuration
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Log      LogConfig      `json:"log"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level string `json:"level"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Debug("No .env file found, using environment variables")
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			Name:     getEnv("DB_NAME", "masihsama"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "127.0.0.1"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

// GetDSN returns database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
