package services_test

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"interview/internal/models"
	"interview/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	args := m.Called(transaction)
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
		tx.ID = 1 // Simulate DB auto-increment
	})

	result, err := service.CreateTransaction(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.True(t, result.Amount.Equal(decimal.NewFromFloat(100.50)))
	assert.Equal(t, "pending", result.Status)
	assert.Equal(t, uint(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransactionError(t *testing.T) {
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
	assert.Equal(t, expectedTx, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactionNotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return((*models.Transaction)(nil), gorm.ErrRecordNotFound)

	result, err := service.GetTransaction(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "transaction not found")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactionError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return((*models.Transaction)(nil), errors.New("database error"))

	result, err := service.GetTransaction(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get transaction")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactions(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	expectedTxs := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
		{ID: 2, UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "success"},
	}

	filters := models.TransactionFilters{
		UserID: 1,
		Status: "pending",
	}

	mockRepo.On("GetAll", filters).Return(expectedTxs, nil)

	result, err := service.GetTransactions(filters)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTxs, result)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactionsError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	filters := models.TransactionFilters{}

	mockRepo.On("GetAll", filters).Return([]models.Transaction{}, errors.New("database error"))

	result, err := service.GetTransactions(filters)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get transactions")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_GetTransactionsInvalidStatusFilter(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	filters := models.TransactionFilters{
		Status: "invalid_status",
	}

	result, err := service.GetTransactions(filters)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid status filter", err.Error())
	// No repo calls should be made
	mockRepo.AssertNotCalled(t, "GetAll")
}

func TestTransactionService_GetTransactionsWithValidStatusFilters(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	validStatuses := []string{"pending", "success", "failed"}

	for _, status := range validStatuses {
		filters := models.TransactionFilters{
			Status: status,
		}

		expectedTxs := []models.Transaction{
			{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: status},
		}

		mockRepo.On("GetAll", filters).Return(expectedTxs, nil).Once()

		result, err := service.GetTransactions(filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedTxs, result)
	}

	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	// Test successful update
	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{ID: 1, Status: "pending"}, nil)
	mockRepo.On("Update", uint(1), map[string]interface{}{"status": "success"}).Return(nil)

	err := service.UpdateTransactionStatus(1, "success")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatusNotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return((*models.Transaction)(nil), gorm.ErrRecordNotFound)

	err := service.UpdateTransactionStatus(1, "success")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction not found")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatusInvalidStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	err := service.UpdateTransactionStatus(1, "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestTransactionService_UpdateTransactionStatusSameStatus(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{ID: 1, Status: "success"}, nil)
	mockRepo.On("Update", uint(1), map[string]interface{}{"status": "success"}).Return(nil)

	err := service.UpdateTransactionStatus(1, "success")

	assert.NoError(t, err) // The service allows setting the same status
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatusGetByIDError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	// Test when GetByID returns an error (not ErrRecordNotFound)
	mockRepo.On("GetByID", uint(1)).Return((*models.Transaction)(nil), errors.New("database connection error"))

	err := service.UpdateTransactionStatus(1, "success")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get transaction")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatusUpdateError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	existingTx := &models.Transaction{
		ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingTx, nil)
	mockRepo.On("Update", uint(1), map[string]interface{}{"status": "success"}).Return(errors.New("update failed"))

	err := service.UpdateTransactionStatus(1, "success")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update transaction")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{ID: 1}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteTransaction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransactionNotFound(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return((*models.Transaction)(nil), gorm.ErrRecordNotFound)

	err := service.DeleteTransaction(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction not found")
	mockRepo.AssertNotCalled(t, "Delete")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransactionError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{ID: 1}, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("database error"))

	err := service.DeleteTransaction(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete transaction")
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_UpdateTransactionStatusCheckExistingTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{
		ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending",
	}, nil)
	mockRepo.On("Update", uint(1), map[string]interface{}{"status": "success"}).Return(nil)

	err := service.UpdateTransactionStatus(1, "success")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransactionCheckExisting(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(&models.Transaction{
		ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending",
	}, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteTransaction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionService_DeleteTransaction_GetByIDError(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("database error"))

	err := service.DeleteTransaction(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get transaction")
	mockRepo.AssertExpectations(t)
}
