package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"

	"interview/internal/handlers"
	"interview/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardService is a mock implementation of DashboardService
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

func setupDashboardTestRouter() (*gin.Engine, *MockDashboardService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockDashboardService)
	handler := handlers.NewDashboardHandler(mockService)

	api := router.Group("/api")
	{
		api.GET("/dashboard/summary", handler.GetSummary)
	}

	return router, mockService
}

func TestDashboardHandler_GetSummary(t *testing.T) {
	router, mockService := setupDashboardTestRouter()

	expectedSummary := &models.DashboardSummary{
		TodaySuccessfulTransactions: 10,
		TodaySuccessfulAmount:       decimal.NewFromFloat(1500.50),
		AverageTransactionPerUser:   decimal.NewFromFloat(2.5),
		LatestTransactions: []models.Transaction{
			{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
			{ID: 2, UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		},
		StatusCounts: models.StatusCounts{
			Success: 5,
			Pending: 3,
			Failed:  2,
		},
	}

	mockService.On("GetSummary").Return(expectedSummary, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/dashboard/summary", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_GetSummaryServiceError(t *testing.T) {
	router, mockService := setupDashboardTestRouter()

	mockService.On("GetSummary").Return((*models.DashboardSummary)(nil), errors.New("service error"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/dashboard/summary", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
