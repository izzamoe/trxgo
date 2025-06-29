package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction represents the transaction model
type Transaction struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id" gorm:"not null;index"`
	Amount    decimal.Decimal `json:"amount" gorm:"not null;type:decimal(15,2)"`
	Status    string          `json:"status" gorm:"not null;default:'pending';index"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// TransactionFilters represents filters for transaction queries
type TransactionFilters struct {
	UserID uint   `form:"user_id"`
	Status string `form:"status"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

// CreateTransactionRequest represents request body for creating transaction
type CreateTransactionRequest struct {
	UserID uint            `json:"user_id" validate:"required,min=1"`
	Amount decimal.Decimal `json:"amount" validate:"required,decimal_positive"`
}

// UpdateTransactionRequest represents request body for updating transaction
type UpdateTransactionRequest struct {
	Status string `json:"status" validate:"required,oneof=pending success failed"`
}

// DashboardSummary represents dashboard summary response
type DashboardSummary struct {
	TodaySuccessfulTransactions int             `json:"today_successful_transactions"`
	TodaySuccessfulAmount       decimal.Decimal `json:"today_successful_amount"`
	AverageTransactionPerUser   decimal.Decimal `json:"average_transaction_per_user"`
	LatestTransactions          []Transaction   `json:"latest_transactions"`
	StatusCounts                StatusCounts    `json:"status_counts"`
}

// StatusCounts represents transaction status counts
type StatusCounts struct {
	Success int `json:"success"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}
