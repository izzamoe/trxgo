package config_test

import (
	"os"
	"testing"

	"interview/internal/config"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear environment variables to test defaults
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("SERVER_PORT")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Database.Host != "127.0.0.1" {
		t.Errorf("Expected default DB host to be '127.0.0.1', got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 3306 {
		t.Errorf("Expected default DB port to be 3306, got %d", cfg.Database.Port)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default server port to be '8080', got %s", cfg.Server.Port)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("SERVER_PORT")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected DB host to be 'localhost', got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 3307 {
		t.Errorf("Expected DB port to be 3307, got %d", cfg.Database.Port)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected server port to be '9090', got %s", cfg.Server.Port)
	}
}

func TestGetDSN(t *testing.T) {
	dbConfig := config.DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	expectedDSN := "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	actualDSN := dbConfig.GetDSN()

	if actualDSN != expectedDSN {
		t.Errorf("Expected DSN to be '%s', got '%s'", expectedDSN, actualDSN)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	os.Setenv("DB_PORT", "invalid_port")
	defer os.Unsetenv("DB_PORT")

	_, err := config.Load()
	if err == nil {
		t.Error("Expected error for invalid port, got nil")
	}
}
