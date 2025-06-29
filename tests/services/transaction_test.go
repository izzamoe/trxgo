package services_test

import (
	"errors"
	"testing"

	"interview/internal/models"
	"interview/internal/services"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repository for testing
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(tx *models.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetAll(filters models.TransactionFilters) ([]models.Transaction, error) {
	args := m.Called(filters)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTodaySuccessful() (int, decimal.Decimal, error) {
	args := m.Called()
	return args.Int(0), args.Get(1).(decimal.Decimal), args.Error(2)
}

func (m *MockTransactionRepository) GetAveragePerUser() (decimal.Decimal, error) {
	args := m.Called()
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockTransactionRepository) GetLatest(limit int) ([]models.Transaction, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetStatusCounts() (models.StatusCounts, error) {
	args := m.Called()
	return args.Get(0).(models.StatusCounts), args.Error(1)
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(nil).Run(func(args mock.Arguments) {
		tx := args.Get(0).(*models.Transaction)
		tx.ID = 1
	})

	result, err := service.CreateTransaction(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.True(t, result.Amount.Equal(decimal.NewFromFloat(100.50)))
	assert.Equal(t, "pending", result.Status)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransaction_Error(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Transaction")).Return(errors.New("database error"))

	result, err := service.CreateTransaction(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create transaction")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	expectedTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockRepo.On("GetByID", uint(1)).Return(expectedTx, nil)

	result, err := service.GetTransaction(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTx.ID, result.ID)
	assert.Equal(t, expectedTx.UserID, result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransaction_NotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetTransaction(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "transaction not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactions(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	filters := models.TransactionFilters{
		UserID: 1,
		Status: "pending",
		Limit:  10,
		Offset: 0,
	}

	expectedTxs := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
		{ID: 2, UserID: 1, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
	}

	mockRepo.On("GetAll", filters).Return(expectedTxs, nil)

	result, err := service.GetTransactions(filters)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedTxs[0].ID, result[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactions_InvalidStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	filters := models.TransactionFilters{
		Status: "invalid_status",
	}

	result, err := service.GetTransactions(filters)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid status filter", err.Error())
}

func TestTransactionService_UpdateTransactionStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	existingTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingTx, nil)
	mockRepo.On("Update", uint(1), map[string]interface{}{"status": "success"}).Return(nil)

	err := service.UpdateTransactionStatus(1, "success")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatus_InvalidStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	err := service.UpdateTransactionStatus(1, "invalid_status")

	assert.Error(t, err)
	assert.Equal(t, "invalid status", err.Error())
}

func TestTransactionService_UpdateTransactionStatus_NotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	err := service.UpdateTransactionStatus(1, "success")

	assert.Error(t, err)
	assert.Equal(t, "transaction not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	existingTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingTx, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteTransaction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransaction_NotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)

	err := service.DeleteTransaction(1)

	assert.Error(t, err)
	assert.Equal(t, "transaction not found", err.Error())
	mockRepo.AssertExpectations(t)
}
