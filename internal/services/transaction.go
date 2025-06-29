package services

import (
	"errors"
	"fmt"

	"interview/internal/models"
	"interview/internal/repositories"

	"gorm.io/gorm"
)

// TransactionService interface defines transaction service methods
type TransactionService interface {
	CreateTransaction(req models.CreateTransactionRequest) (*models.Transaction, error)
	GetTransaction(id uint) (*models.Transaction, error)
	GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error)
	UpdateTransactionStatus(id uint, status string) error
	DeleteTransaction(id uint) error
}

// transactionService implements TransactionService interface
type transactionService struct {
	repo repositories.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// CreateTransaction creates a new transaction
func (s *transactionService) CreateTransaction(req models.CreateTransactionRequest) (*models.Transaction, error) {
	transaction := &models.Transaction{
		UserID: req.UserID,
		Amount: req.Amount,
		Status: "pending",
	}

	err := s.repo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return transaction, nil
}

// GetTransaction gets a transaction by ID
func (s *transactionService) GetTransaction(id uint) (*models.Transaction, error) {
	transaction, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	return transaction, nil
}

// GetTransactions gets all transactions with filters
func (s *transactionService) GetTransactions(filters models.TransactionFilters) ([]models.Transaction, error) {
	// Validate status filter
	if filters.Status != "" && filters.Status != "pending" && filters.Status != "success" && filters.Status != "failed" {
		return nil, errors.New("invalid status filter")
	}

	transactions, err := s.repo.GetAll(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}

	return transactions, nil
}

// UpdateTransactionStatus updates transaction status
func (s *transactionService) UpdateTransactionStatus(id uint, status string) error {
	// Validate status
	if status != "pending" && status != "success" && status != "failed" {
		return errors.New("invalid status")
	}

	// Check if transaction exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	updates := map[string]interface{}{
		"status": status,
	}

	err = s.repo.Update(id, updates)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	return nil
}

// DeleteTransaction deletes a transaction
func (s *transactionService) DeleteTransaction(id uint) error {
	// Check if transaction exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	return nil
}
