package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"interview/internal/config"
	"interview/internal/handlers"
	"interview/internal/middleware"
	"interview/internal/models"
	"interview/internal/repositories"
	"interview/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup logging
	setupLogging(cfg.Log.Level)

	// Initialize database
	db, err := initializeDatabase(cfg.Database)
	if err != nil {
		logrus.Fatal("Failed to initialize database:", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		logrus.Fatal("Failed to migrate database:", err)
	}

	// Initialize dependencies
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	dashboardService := services.NewDashboardService(transactionRepo)

	transactionHandler := handlers.NewTransactionHandler(transactionService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Setup router
	router := setupRouter(transactionHandler, dashboardHandler)

	// Start server
	address := cfg.Server.Host + ":" + cfg.Server.Port
	logrus.Info("Starting server on ", address)
	if err := router.Run(address); err != nil {
		logrus.Fatal("Failed to start server:", err)
	}
}

// setupLogging configures the logging system
func setupLogging(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
}

// initializeDatabase initializes the database connection
func initializeDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	logrus.Info("Database connection established")
	return db, nil
}

// setupRouter configures the HTTP router
func setupRouter(transactionHandler *handlers.TransactionHandler, dashboardHandler *handlers.DashboardHandler) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// API routes
	api := router.Group("/api")
	{
		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/summary", dashboardHandler.GetSummary)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	return router
}
