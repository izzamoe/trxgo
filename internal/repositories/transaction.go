package repositories

import (
	"time"

	"interview/internal/models"

	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

// TransactionRepository interface defines transaction repository methods
type TransactionRepository interface {
	Create(tx *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	GetAll(filters models.TransactionFilters) ([]models.Transaction, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	GetTodaySuccessful() (int, decimal.Decimal, error)
	GetAveragePerUser() (decimal.Decimal, error)
	GetLatest(limit int) ([]models.Transaction, error)
	GetStatusCounts() (models.StatusCounts, error)
}

// transactionRepository implements TransactionRepository interface
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create creates a new transaction
func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

// GetByID gets a transaction by ID
func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAll gets all transactions with filters
func (r *transactionRepository) GetAll(filters models.TransactionFilters) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Model(&models.Transaction{})

	// Apply filters
	if filters.UserID != 0 {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Set default limit and offset
	limit := filters.Limit
	if limit == 0 || limit > 100 {
		limit = 20
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

// Update updates a transaction
func (r *transactionRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a transaction
func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// GetTodaySuccessful gets today's successful transactions count and amount
func (r *transactionRepository) GetTodaySuccessful() (int, decimal.Decimal, error) {
	var count int64
	var totalAmount decimal.Decimal

	today := time.Now().Format("2006-01-02")

	err := r.db.Model(&models.Transaction{}).
		Where("status = ? AND DATE(created_at) = ?", "success", today).
		Count(&count).Error
	if err != nil {
		return 0, decimal.Zero, err
	}

	err = r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ? AND DATE(created_at) = ?", "success", today).
		Scan(&totalAmount).Error
	if err != nil {
		return 0, decimal.Zero, err
	}

	return int(count), totalAmount, nil
}

// GetAveragePerUser gets average transactions per user
func (r *transactionRepository) GetAveragePerUser() (decimal.Decimal, error) {
	var result struct {
		Average decimal.Decimal
	}

	err := r.db.Raw(`
		SELECT COALESCE(AVG(user_transaction_count), 0) as average
		FROM (
			SELECT user_id, COUNT(*) as user_transaction_count
			FROM transactions
			GROUP BY user_id
		) as user_counts
	`).Scan(&result).Error

	return result.Average, err
}

// GetLatest gets latest transactions
func (r *transactionRepository) GetLatest(limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Order("created_at DESC").Limit(limit).Find(&transactions).Error
	return transactions, err
}

// GetStatusCounts gets transaction counts by status
func (r *transactionRepository) GetStatusCounts() (models.StatusCounts, error) {
	var counts models.StatusCounts

	err := r.db.Model(&models.Transaction{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&[]struct {
			Status string
			Count  int
		}{}).Error

	if err != nil {
		return counts, err
	}

	// Get individual counts
	var successCount, pendingCount, failedCount int64
	r.db.Model(&models.Transaction{}).Where("status = ?", "success").Count(&successCount)
	r.db.Model(&models.Transaction{}).Where("status = ?", "pending").Count(&pendingCount)
	r.db.Model(&models.Transaction{}).Where("status = ?", "failed").Count(&failedCount)

	counts.Success = int(successCount)
	counts.Pending = int(pendingCount)
	counts.Failed = int(failedCount)

	return counts, nil
}
