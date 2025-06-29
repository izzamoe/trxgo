package services_test

import (
	"errors"
	"testing"

	"interview/internal/models"
	"interview/internal/services"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestDashboardService_GetSummary(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewDashboardService(mockRepo)

	expectedTransactions := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{ID: 2, UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
	}

	expectedStatusCounts := models.StatusCounts{
		Success: 5,
		Pending: 3,
		Failed:  2,
	}

	mockRepo.On("GetTodaySuccessful").Return(10, decimal.NewFromFloat(1500.50), nil)
	mockRepo.On("GetAveragePerUser").Return(decimal.NewFromFloat(2.5), nil)
	mockRepo.On("GetLatest", 10).Return(expectedTransactions, nil)
	mockRepo.On("GetStatusCounts").Return(expectedStatusCounts, nil)

	result, err := service.GetSummary()

	assert.NoError(t, err)
	assert.Equal(t, 10, result.TodaySuccessfulTransactions)
	assert.True(t, result.TodaySuccessfulAmount.Equal(decimal.NewFromFloat(1500.50)))
	assert.True(t, result.AverageTransactionPerUser.Equal(decimal.NewFromFloat(2.5)))
	assert.Equal(t, expectedTransactions, result.LatestTransactions)
	assert.Equal(t, expectedStatusCounts, result.StatusCounts)
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummaryTodaySuccessfulError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewDashboardService(mockRepo)

	mockRepo.On("GetTodaySuccessful").Return(0, decimal.Zero, errors.New("database error"))

	result, err := service.GetSummary()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get today's successful transactions")
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummaryAveragePerUserError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewDashboardService(mockRepo)

	mockRepo.On("GetTodaySuccessful").Return(5, decimal.NewFromFloat(500.00), nil)
	mockRepo.On("GetAveragePerUser").Return(decimal.Zero, errors.New("calculation error"))

	result, err := service.GetSummary()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get average transactions per user")
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummaryLatestTransactionsError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewDashboardService(mockRepo)

	mockRepo.On("GetTodaySuccessful").Return(5, decimal.NewFromFloat(500.00), nil)
	mockRepo.On("GetAveragePerUser").Return(decimal.NewFromFloat(3.0), nil)
	mockRepo.On("GetLatest", 10).Return([]models.Transaction{}, errors.New("fetch error"))

	result, err := service.GetSummary()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get latest transactions")
	mockRepo.AssertExpectations(t)
}

func TestDashboardService_GetSummaryStatusCountsError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewDashboardService(mockRepo)

	expectedTransactions := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
	}

	mockRepo.On("GetTodaySuccessful").Return(5, decimal.NewFromFloat(500.00), nil)
	mockRepo.On("GetAveragePerUser").Return(decimal.NewFromFloat(3.0), nil)
	mockRepo.On("GetLatest", 10).Return(expectedTransactions, nil)
	mockRepo.On("GetStatusCounts").Return(models.StatusCounts{}, errors.New("count error"))

	result, err := service.GetSummary()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get status counts")
	mockRepo.AssertExpectations(t)
}
