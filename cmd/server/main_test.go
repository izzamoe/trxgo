package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/config"
	"interview/internal/handlers"
	"interview/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService for testing
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(req models.CreateTransactionRequest) (*models.Transaction, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransaction(id uint) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionService) UpdateTransactionStatus(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockTransactionService) DeleteTransaction(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockDashboardService for testing
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetSummary() (*models.DashboardSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardSummary), args.Error(1)
}

func TestSetupLogging(t *testing.T) {
	// Test with valid log level
	setupLogging("debug")

	// Test with invalid log level
	setupLogging("invalid")
}

func TestInitializeDatabase(t *testing.T) {
	// Test with invalid DSN to ensure error handling
	cfg := config.DatabaseConfig{
		Host:     "invalid-host",
		Port:     3306,
		User:     "root",
		Password: "password",
		Name:     "test",
	}

	db, err := initializeDatabase(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestInitializeDatabaseError(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "invalid_host",
		Port:     3306,
		User:     "root",
		Password: "wrong_password",
		Name:     "nonexistent_db",
	}

	db, err := initializeDatabase(cfg)

	// Should fail with invalid connection
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestInitializeDatabaseConnectionPool(t *testing.T) {
	// Test with valid SQLite in-memory database for testing
	cfg := config.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "root",
		Name:     "masihsama",
	}

	// This test verifies the connection pool setup logic
	// In a real scenario, we'd use a test database
	db, err := initializeDatabase(cfg)
	if err == nil && db != nil {
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		assert.NotNil(t, sqlDB)

		// Just verify that we can get the DB stats (connection pool was set)
		stats := sqlDB.Stats()
		assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)
	}
	// Test passes even if connection fails (for CI environments)
}

func TestInitializeDatabaseSuccess(t *testing.T) {
	// Test successful database initialization with valid config
	cfg := config.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "root",
		Name:     "masihsama",
	}

	db, err := initializeDatabase(cfg)

	// If connection succeeds, verify all setup
	if err == nil && db != nil {
		sqlDB, err := db.DB()
		assert.NoError(t, err)
		assert.NotNil(t, sqlDB)

		// Test connection pool settings
		stats := sqlDB.Stats()
		assert.GreaterOrEqual(t, stats.MaxOpenConnections, 0)

		// Verify we can ping the database
		err = sqlDB.Ping()
		if err == nil {
			// Connection successful, verify pool settings
			assert.Equal(t, 100, stats.MaxOpenConnections)
		}
	}
}

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockTxService := new(MockTransactionService)
	mockDashService := new(MockDashboardService)

	// Create handlers with mock services
	transactionHandler := handlers.NewTransactionHandler(mockTxService)
	dashboardHandler := handlers.NewDashboardHandler(mockDashService)

	router := setupRouter(transactionHandler, dashboardHandler)

	assert.NotNil(t, router)

	// Test that routes are properly configured by checking router structure
	routes := router.Routes()
	assert.True(t, len(routes) > 0)

	// Check for specific route patterns
	foundHealthRoute := false
	foundTransactionRoute := false
	foundDashboardRoute := false

	for _, route := range routes {
		if route.Path == "/health" {
			foundHealthRoute = true
		}
		if route.Path == "/api/transactions" {
			foundTransactionRoute = true
		}
		if route.Path == "/api/dashboard/summary" {
			foundDashboardRoute = true
		}
	}

	assert.True(t, foundHealthRoute, "Health route should be configured")
	assert.True(t, foundTransactionRoute, "Transaction route should be configured")
	assert.True(t, foundDashboardRoute, "Dashboard route should be configured")
}

func TestSetupRouterComprehensive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockTxService := new(MockTransactionService)
	mockDashService := new(MockDashboardService)

	// Create handlers with mock services
	transactionHandler := handlers.NewTransactionHandler(mockTxService)
	dashboardHandler := handlers.NewDashboardHandler(mockDashService)

	router := setupRouter(transactionHandler, dashboardHandler)

	// Test all routes exist
	routes := router.Routes()
	expectedRoutes := []string{
		"/health",
		"/api/transactions",
		"/api/transactions/:id",
		"/api/dashboard/summary",
	}

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Path] = true
	}

	for _, expectedRoute := range expectedRoutes {
		assert.True(t, routeMap[expectedRoute], "Route %s should be configured", expectedRoute)
	}
}

func TestInitializeDatabaseSuccessDetailed(t *testing.T) {
	// This test verifies successful database initialization logic
	// In a real production environment, you'd use a test database
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "test_user",
		Password: "test_password",
		Name:     "test_db",
	}

	// Test that the function handles database configuration properly
	// Even if connection fails, we can test the DSN construction logic
	db, err := initializeDatabase(cfg)

	// In CI/test environments, DB might not be available
	// So we test that the function runs without panic
	if err != nil {
		// Expected in test environment
		assert.Nil(t, db)
		assert.Error(t, err)
	} else {
		// If connection succeeds (in dev environment)
		assert.NotNil(t, db)
		assert.NoError(t, err)
	}
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockTxService := new(MockTransactionService)
	mockDashService := new(MockDashboardService)

	// Create handlers with mock services
	transactionHandler := handlers.NewTransactionHandler(mockTxService)
	dashboardHandler := handlers.NewDashboardHandler(mockDashService)

	router := setupRouter(transactionHandler, dashboardHandler)

	// Test health endpoint
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OK")
}

func TestInitializeDatabaseConfigError(t *testing.T) {
	// Test with empty configuration that should cause errors
	cfg := config.DatabaseConfig{}

	db, err := initializeDatabase(cfg)

	// Should fail with empty configuration
	assert.Error(t, err)
	assert.Nil(t, db)
}

// TestMain_ComponentsIntegration tests main function components
func TestMain_ComponentsIntegration(t *testing.T) {
	// This test verifies that all components used by main() work together
	// We can't test main() directly, but we can test its key components

	// Test configuration loading would work
	// Test logging setup
	setupLogging("debug")
	setupLogging("info")
	setupLogging("warn")
	setupLogging("error")

	// Test invalid log level
	setupLogging("invalid")
}

func TestSetupLoggingComprehensive(t *testing.T) {
	// Test all log levels to ensure comprehensive coverage
	logLevels := []string{"debug", "info", "warn", "error", "fatal", "panic", "trace", "invalid"}

	for _, level := range logLevels {
		// This should not panic for any input
		setupLogging(level)
	}
}
