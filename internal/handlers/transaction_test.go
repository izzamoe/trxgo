package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"interview/internal/handlers"
	"interview/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService is a mock implementation of TransactionService
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

func setupTestRouter() (*gin.Engine, *MockTransactionService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	api := router.Group("/api")
	{
		api.POST("/transactions", handler.CreateTransaction)
		api.GET("/transactions", handler.GetTransactions)
		api.GET("/transactions/:id", handler.GetTransaction)
		api.PUT("/transactions/:id", handler.UpdateTransaction)
		api.DELETE("/transactions/:id", handler.DeleteTransaction)
	}

	return router, mockService
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	router, mockService := setupTestRouter()

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

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_CreateTransactionInvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	body := bytes.NewBufferString("invalid json")
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/transactions", body)
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_CreateTransactionValidationFailed(t *testing.T) {
	router, _ := setupTestRouter()

	req := models.CreateTransactionRequest{
		UserID: 0, // Invalid - should be >= 1
		Amount: decimal.NewFromFloat(100.50),
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_CreateTransactionServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
	}

	mockService.On("CreateTransaction", req).Return((*models.Transaction)(nil), errors.New("service error"))

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransaction(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedTx := &models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	mockService.On("GetTransaction", uint(1)).Return(expectedTx, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactionInvalidID(t *testing.T) {
	router, _ := setupTestRouter()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions/invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_GetTransactionNotFound(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("GetTransaction", uint(1)).Return((*models.Transaction)(nil), errors.New("transaction not found"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactionServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("GetTransaction", uint(1)).Return((*models.Transaction)(nil), errors.New("service error"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedTxs := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
		{ID: 2, UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "success"},
	}

	mockService.On("GetTransactions", mock.AnythingOfType("models.TransactionFilters")).Return(expectedTxs, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactionsWithFilters(t *testing.T) {
	router, mockService := setupTestRouter()

	expectedTxs := []models.Transaction{
		{ID: 1, UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
	}

	mockService.On("GetTransactions", mock.AnythingOfType("models.TransactionFilters")).Return(expectedTxs, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions?user_id=1&status=pending&limit=10&offset=0", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_GetTransactionsInvalidQuery(t *testing.T) {
	router, _ := setupTestRouter()

	// Test with invalid query parameters
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions?user_id=invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_GetTransactionsServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("GetTransactions", mock.AnythingOfType("models.TransactionFilters")).Return([]models.Transaction{}, errors.New("service error"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/transactions", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code) // Handler returns BadRequest for service errors
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransaction(t *testing.T) {
	router, mockService := setupTestRouter()

	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransactionInvalidID(t *testing.T) {
	router, _ := setupTestRouter()

	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/invalid", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_UpdateTransactionInvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	body := bytes.NewBufferString("invalid json")
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/1", body)
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_UpdateTransactionValidationFailed(t *testing.T) {
	router, _ := setupTestRouter()

	req := models.UpdateTransactionRequest{
		Status: "invalid", // Invalid status
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_UpdateTransactionNotFound(t *testing.T) {
	router, mockService := setupTestRouter()

	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(errors.New("transaction not found"))

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransactionServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(errors.New("service error"))

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/transactions/1", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("DeleteTransaction", uint(1)).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_DeleteTransactionInvalidID(t *testing.T) {
	router, _ := setupTestRouter()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/transactions/invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_DeleteTransactionNotFound(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("DeleteTransaction", uint(1)).Return(errors.New("transaction not found"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_DeleteTransactionServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("DeleteTransaction", uint(1)).Return(errors.New("service error"))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/transactions/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransactionInvalidStatus(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router, _ := setupTestRouter()
	router.PUT("/transactions/:id", handler.UpdateTransaction)

	// Don't set expectation since validation should fail before service call

	reqBody := `{"status": "invalid_status"}`
	req, _ := http.NewRequest("PUT", "/transactions/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Don't assert expectations since service should not be called
}

func TestTransactionHandler_UpdateTransactionServiceNotFound(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router, _ := setupTestRouter()
	router.PUT("/transactions/:id", handler.UpdateTransaction)

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(errors.New("transaction not found"))

	reqBody := `{"status": "success"}`
	req, _ := http.NewRequest("PUT", "/transactions/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransactionInvalidStatusFromService(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router, _ := setupTestRouter()
	router.PUT("/transactions/:id", handler.UpdateTransaction)

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(errors.New("invalid status"))

	reqBody := `{"status": "success"}`
	req, _ := http.NewRequest("PUT", "/transactions/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestTransactionHandler_UpdateTransactionInternalError(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	router, _ := setupTestRouter()
	router.PUT("/transactions/:id", handler.UpdateTransaction)

	mockService.On("UpdateTransactionStatus", uint(1), "success").Return(errors.New("database error"))

	reqBody := `{"status": "success"}`
	req, _ := http.NewRequest("PUT", "/transactions/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
