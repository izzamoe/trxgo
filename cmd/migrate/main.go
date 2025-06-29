package main

import (
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"interview/internal/config"
	"interview/internal/models"
)

func main() {
	var action string
	var verbose bool
	flag.StringVar(&action, "action", "up", "Migration action: up, down, reset, status")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := connectDB(cfg.Database, verbose)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch action {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("✅ Migrations completed successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("✅ Migrations rolled back successfully")
	case "reset":
		if err := migrateReset(db); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		fmt.Println("✅ Database reset successfully")
	case "status":
		if err := migrateStatus(db); err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}
	default:
		log.Fatalf("Unknown action: %s. Use: up, down, reset, or status", action)
	}
}

func connectDB(cfg config.DatabaseConfig, verbose bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	logMode := logger.Silent
	if verbose {
		logMode = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func migrateUp(db *gorm.DB) error {
	fmt.Println("🚀 Running migrations...")

	// Auto migrate all models
	if err := db.AutoMigrate(
		&models.Transaction{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create additional indexes if needed
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func migrateDown(db *gorm.DB) error {
	fmt.Println("📉 Rolling back migrations...")

	// Drop all tables in reverse order
	if err := db.Migrator().DropTable(&models.Transaction{}); err != nil {
		return fmt.Errorf("failed to drop transactions table: %w", err)
	}

	return nil
}

func migrateReset(db *gorm.DB) error {
	fmt.Println("🔄 Resetting database...")

	// Drop all tables
	if err := migrateDown(db); err != nil {
		return err
	}

	// Run migrations again
	if err := migrateUp(db); err != nil {
		return err
	}

	return nil
}

func createIndexes(db *gorm.DB) error {
	// Check if we need to create any additional indexes
	// GORM should handle most indexes from struct tags, but we can add custom ones here

	// Example: Create composite index
	if !db.Migrator().HasIndex(&models.Transaction{}, "idx_user_status") {
		if err := db.Exec("CREATE INDEX idx_user_status ON transactions (user_id, status)").Error; err != nil {
			return fmt.Errorf("failed to create composite index: %w", err)
		}
	}

	return nil
}

func migrateStatus(db *gorm.DB) error {
	fmt.Println("📊 Migration Status")
	fmt.Println("==================")

	// Check if tables exist
	tables := []interface{}{
		&models.Transaction{},
	}

	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			fmt.Printf("✅ Table exists: %T\n", table)

			// Get column info
			columnTypes, err := db.Migrator().ColumnTypes(table)
			if err == nil {
				fmt.Printf("   Columns: %d\n", len(columnTypes))
			}

			// Check indexes
			indexes, err := db.Migrator().GetIndexes(table)
			if err == nil {
				fmt.Printf("   Indexes: %d\n", len(indexes))
			}
		} else {
			fmt.Printf("❌ Table missing: %T\n", table)
		}
	}

	fmt.Println("==================")
	return nil
}
