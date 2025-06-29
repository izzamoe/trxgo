package models_test

import (
	"testing"

	"interview/internal/models"

	"github.com/shopspring/decimal"
)

func TestTransactionModel(t *testing.T) {
	amount := decimal.NewFromFloat(100.50)
	transaction := models.Transaction{
		ID:     1,
		UserID: 1,
		Amount: amount,
		Status: "pending",
	}

	if transaction.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", transaction.ID)
	}
	if transaction.UserID != 1 {
		t.Errorf("Expected UserID to be 1, got %d", transaction.UserID)
	}
	if !transaction.Amount.Equal(amount) {
		t.Errorf("Expected Amount to be %s, got %s", amount.String(), transaction.Amount.String())
	}
	if transaction.Status != "pending" {
		t.Errorf("Expected Status to be 'pending', got %s", transaction.Status)
	}
}

func TestAPIResponse(t *testing.T) {
	response := models.APIResponse{
		Success: true,
		Data:    "test data",
		Message: "test message",
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}
	if response.Data != "test data" {
		t.Errorf("Expected Data to be 'test data', got %v", response.Data)
	}
	if response.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %s", response.Message)
	}
}

func TestDashboardSummary(t *testing.T) {
	todayAmount := decimal.NewFromFloat(1250.75)
	avgPerUser := decimal.NewFromFloat(3.2)
	summary := models.DashboardSummary{
		TodaySuccessfulTransactions: 5,
		TodaySuccessfulAmount:       todayAmount,
		AverageTransactionPerUser:   avgPerUser,
		StatusCounts: models.StatusCounts{
			Success: 15,
			Pending: 8,
			Failed:  2,
		},
	}

	if summary.TodaySuccessfulTransactions != 5 {
		t.Errorf("Expected TodaySuccessfulTransactions to be 5, got %d", summary.TodaySuccessfulTransactions)
	}
	if !summary.TodaySuccessfulAmount.Equal(todayAmount) {
		t.Errorf("Expected TodaySuccessfulAmount to be %s, got %s", todayAmount.String(), summary.TodaySuccessfulAmount.String())
	}
	if summary.StatusCounts.Success != 15 {
		t.Errorf("Expected Success count to be 15, got %d", summary.StatusCounts.Success)
	}
}

func TestTransactionFilters(t *testing.T) {
	filters := models.TransactionFilters{
		UserID: 1,
		Status: "pending",
		Limit:  20,
		Offset: 0,
	}

	if filters.UserID != 1 {
		t.Errorf("Expected UserID to be 1, got %d", filters.UserID)
	}
	if filters.Status != "pending" {
		t.Errorf("Expected Status to be 'pending', got %s", filters.Status)
	}
}

func TestCreateTransactionRequest(t *testing.T) {
	amount := decimal.NewFromFloat(100.50)
	req := models.CreateTransactionRequest{
		UserID: 1,
		Amount: amount,
	}

	if req.UserID != 1 {
		t.Errorf("Expected UserID to be 1, got %d", req.UserID)
	}
	if !req.Amount.Equal(amount) {
		t.Errorf("Expected Amount to be %s, got %s", amount.String(), req.Amount.String())
	}
}

func TestUpdateTransactionRequest(t *testing.T) {
	req := models.UpdateTransactionRequest{
		Status: "success",
	}

	if req.Status != "success" {
		t.Errorf("Expected Status to be 'success', got %s", req.Status)
	}
}

func TestStatusCounts(t *testing.T) {
	counts := models.StatusCounts{
		Success: 10,
		Pending: 5,
		Failed:  2,
	}

	total := counts.Success + counts.Pending + counts.Failed
	if total != 17 {
		t.Errorf("Expected total count to be 17, got %d", total)
	}
}
