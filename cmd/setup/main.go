package main

import (
	"fmt"
	"log"

	"interview/internal/config"
	"interview/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🚀 Database Setup Tool")
	fmt.Println("======================")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect without database name first to create database
	if err := createDatabase(cfg.Database); err != nil {
		log.Printf("Warning: Could not create database: %v", err)
		fmt.Println("📝 Please create the database manually:")
		fmt.Printf("   mysql -u %s -p -e 'CREATE DATABASE IF NOT EXISTS %s;'\n", cfg.Database.User, cfg.Database.Name)
	}

	// Connect to the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Connected to database")

	// Run migrations
	fmt.Println("🏗️  Running migrations...")
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	fmt.Println("✅ Database setup completed successfully!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'make run' to start the server")
	fmt.Println("  2. Visit http://localhost:8080/health to test")
}

func createDatabase(cfg config.DatabaseConfig) error {
	// Connect without database name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	// Create database
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.Name)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	fmt.Printf("✅ Database '%s' created/verified\n", cfg.Name)
	return nil
}
