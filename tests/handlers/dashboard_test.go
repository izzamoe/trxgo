package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/handlers"
	"interview/internal/models"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dashboard service for testing
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

func TestDashboardHandler_GetSummary(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := handlers.NewDashboardHandler(mockService)

	router := setupTestRouter()
	router.GET("/dashboard/summary", handler.GetSummary)

	expectedSummary := &models.DashboardSummary{
		TodaySuccessfulTransactions: 5,
		TodaySuccessfulAmount:       decimal.NewFromFloat(1250.75),
		AverageTransactionPerUser:   decimal.NewFromFloat(3.2),
		LatestTransactions: []models.Transaction{
			{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		},
		StatusCounts: models.StatusCounts{
			Success: 15,
			Pending: 8,
			Failed:  2,
		},
	}

	mockService.On("GetSummary").Return(expectedSummary, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/dashboard/summary", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Dashboard summary retrieved successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_GetSummary_ServiceError(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := handlers.NewDashboardHandler(mockService)

	router := setupTestRouter()
	router.GET("/dashboard/summary", handler.GetSummary)

	mockService.On("GetSummary").Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/dashboard/summary", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "service error", response.Error)
	mockService.AssertExpectations(t)
}
