package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/handlers"
	"interview/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/shopspring/decimal"
)

// Mock service for testing
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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.POST("/transactions", handler.CreateTransaction)

	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
	}

	expectedTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockService.On("CreateTransaction", req).Return(expectedTx, nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Transaction created successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransaction_InvalidJSON(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.POST("/transactions", handler.CreateTransaction)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid request body", response.Error)
}

func TestTransactionHandler_CreateTransaction_ValidationError(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.POST("/transactions", handler.CreateTransaction)

	req := models.CreateTransactionRequest{
		UserID: 0, // Invalid - should be min=1
		Amount: decimal.NewFromFloat(100.50),
	}

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "Validation failed")
}

func TestTransactionHandler_CreateTransaction_ServiceError(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.POST("/transactions", handler.CreateTransaction)

	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
	}

	mockService.On("CreateTransaction", req).Return(nil, errors.New("service error"))

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.GET("/transactions", handler.GetTransactions)

	expectedTxs := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
		{ID: 2, UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "success"},
	}

	mockService.On("GetTransactions", mock.AnythingOfType("models.TransactionFilters")).Return(expectedTxs, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/transactions", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Transactions retrieved successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransaction(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.GET("/transactions/:id", handler.GetTransaction)

	expectedTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockService.On("GetTransaction", uint(1)).Return(expectedTx, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Transaction retrieved successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransaction_InvalidID(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.GET("/transactions/:id", handler.GetTransaction)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/transactions/invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid transaction ID", response.Error)
}

func TestTransactionHandler_GetTransaction_NotFound(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.GET("/transactions/:id", handler.GetTransaction)

	mockService.On("GetTransaction", uint(1)).Return(nil, errors.New("transaction not found"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Transaction not found", response.Error)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransaction(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.PUT("/transactions/:id", handler.UpdateTransaction)

	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/transactions/1", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Transaction updated successfully", response.Message)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router := setupTestRouter()
	router.DELETE("/transactions/:id", handler.DeleteTransaction)

	mockService.On("DeleteTransaction", uint(1)).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Transaction deleted successfully", response.Message)
	mockService.AssertExpectations(t)
}
