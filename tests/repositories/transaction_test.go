package repositories_test

import (
	"testing"

	"interview/internal/models"
	"interview/internal/repositories"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repositories.TransactionRepository
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	var err error

	// Use test database connection
	dsn := "root:root@tcp(127.0.0.1:3306)/masihsama?charset=utf8mb4&parseTime=True&loc=Local"
	suite.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(suite.T(), err)

	// Migrate the schema
	err = suite.db.AutoMigrate(&models.Transaction{})
	assert.NoError(suite.T(), err)

	suite.repo = repositories.NewTransactionRepository(suite.db)

	// Clean up any existing test data
	suite.db.Exec("DELETE FROM transactions")
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	// Clean up test data
	suite.db.Exec("DELETE FROM transactions")
}

func (suite *TransactionRepositoryTestSuite) TestCreate() {
	transaction := &models.Transaction{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}

	err := suite.repo.Create(transaction)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), transaction.ID)
}

func (suite *TransactionRepositoryTestSuite) TestGetByID() {
	// Create a test transaction
	transaction := &models.Transaction{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}
	err := suite.repo.Create(transaction)
	assert.NoError(suite.T(), err)

	// Get the transaction
	result, err := suite.repo.GetByID(transaction.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), transaction.ID, result.ID)
	assert.Equal(suite.T(), transaction.UserID, result.UserID)
	assert.True(suite.T(), transaction.Amount.Equal(result.Amount))
}

func (suite *TransactionRepositoryTestSuite) TestGetAll() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "pending"},
		{UserID: 1, Amount: decimal.NewFromFloat(200.00), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(150.75), Status: "pending"},
	}

	for _, tx := range transactions {
		err := suite.repo.Create(tx)
		assert.NoError(suite.T(), err)
	}

	// Test getting all transactions
	filters := models.TransactionFilters{Limit: 10}
	results, err := suite.repo.GetAll(filters)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)

	// Test filtering by user ID
	filters = models.TransactionFilters{UserID: 1, Limit: 10}
	results, err = suite.repo.GetAll(filters)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)

	// Test filtering by status
	filters = models.TransactionFilters{Status: "pending", Limit: 10}
	results, err = suite.repo.GetAll(filters)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
}

func (suite *TransactionRepositoryTestSuite) TestUpdate() {
	// Create a test transaction
	transaction := &models.Transaction{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}
	err := suite.repo.Create(transaction)
	assert.NoError(suite.T(), err)

	// Update the transaction
	updates := map[string]interface{}{
		"status": "success",
	}
	err = suite.repo.Update(transaction.ID, updates)
	assert.NoError(suite.T(), err)

	// Verify the update
	result, err := suite.repo.GetByID(transaction.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "success", result.Status)
}

func (suite *TransactionRepositoryTestSuite) TestDelete() {
	// Create a test transaction
	transaction := &models.Transaction{
		UserID: 1,
		Amount: decimal.NewFromFloat(100.50),
		Status: "pending",
	}
	err := suite.repo.Create(transaction)
	assert.NoError(suite.T(), err)

	// Delete the transaction
	err = suite.repo.Delete(transaction.ID)
	assert.NoError(suite.T(), err)

	// Verify the deletion
	_, err = suite.repo.GetByID(transaction.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *TransactionRepositoryTestSuite) TestGetTodaySuccessful() {
	// Create test transactions for today
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "success"},
		{UserID: 3, Amount: decimal.NewFromFloat(150.75), Status: "pending"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	count, amount, err := suite.repo.GetTodaySuccessful()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
	assert.True(suite.T(), decimal.NewFromFloat(300.50).Equal(amount))
}

func (suite *TransactionRepositoryTestSuite) TestGetTodaySuccessfulNoTransactions() {
	count, amount, err := suite.repo.GetTodaySuccessful()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
	assert.True(suite.T(), decimal.Zero.Equal(amount))
}

func (suite *TransactionRepositoryTestSuite) TestGetAveragePerUser() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 1, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		{UserID: 2, Amount: decimal.NewFromFloat(150.75), Status: "success"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	average, err := suite.repo.GetAveragePerUser()

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), decimal.NewFromFloat(1.5).Equal(average)) // User 1 has 2, User 2 has 1, avg = 1.5
}

func (suite *TransactionRepositoryTestSuite) TestGetAveragePerUserNoTransactions() {
	average, err := suite.repo.GetAveragePerUser()

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), decimal.Zero.Equal(average))
}

func (suite *TransactionRepositoryTestSuite) TestGetLatest() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		{UserID: 3, Amount: decimal.NewFromFloat(150.75), Status: "failed"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	latest, err := suite.repo.GetLatest(2)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), latest, 2)
	// Should be ordered by created_at DESC
	assert.Equal(suite.T(), uint(3), latest[0].UserID)
	assert.Equal(suite.T(), uint(2), latest[1].UserID)
}

func (suite *TransactionRepositoryTestSuite) TestGetLatestEmpty() {
	latest, err := suite.repo.GetLatest(5)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), latest, 0)
}

func (suite *TransactionRepositoryTestSuite) TestGetStatusCounts() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "success"},
		{UserID: 3, Amount: decimal.NewFromFloat(150.75), Status: "pending"},
		{UserID: 4, Amount: decimal.NewFromFloat(75.25), Status: "failed"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	counts, err := suite.repo.GetStatusCounts()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, counts.Success)
	assert.Equal(suite.T(), 1, counts.Pending)
	assert.Equal(suite.T(), 1, counts.Failed)
}

func (suite *TransactionRepositoryTestSuite) TestGetStatusCountsEmpty() {
	counts, err := suite.repo.GetStatusCounts()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, counts.Success)
	assert.Equal(suite.T(), 0, counts.Pending)
	assert.Equal(suite.T(), 0, counts.Failed)
}

// Test GetAll edge cases for better coverage
func (suite *TransactionRepositoryTestSuite) TestGetAllWithLimitAndOffset() {
	// Create more test transactions
	for i := 1; i <= 5; i++ {
		tx := &models.Transaction{
			UserID: uint(i),
			Amount: decimal.NewFromInt(int64(i * 100)),
			Status: "success",
		}
		suite.repo.Create(tx)
	}

	// Test with limit and offset
	filters := models.TransactionFilters{
		Limit:  2,
		Offset: 2,
	}

	transactions, err := suite.repo.GetAll(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), transactions, 2)
}

func (suite *TransactionRepositoryTestSuite) TestGetAllWithUserIDFilter() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 1, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		{UserID: 2, Amount: decimal.NewFromFloat(150.75), Status: "success"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// Test with UserID filter
	filters := models.TransactionFilters{
		UserID: 1,
	}

	result, err := suite.repo.GetAll(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	for _, tx := range result {
		assert.Equal(suite.T(), uint(1), tx.UserID)
	}
}

func (suite *TransactionRepositoryTestSuite) TestGetAllWithStatusFilter() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		{UserID: 3, Amount: decimal.NewFromFloat(150.75), Status: "success"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// Test with Status filter
	filters := models.TransactionFilters{
		Status: "success",
	}

	result, err := suite.repo.GetAll(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	for _, tx := range result {
		assert.Equal(suite.T(), "success", tx.Status)
	}
}

func (suite *TransactionRepositoryTestSuite) TestGetAllWithAllFilters() {
	// Create test transactions
	transactions := []*models.Transaction{
		{UserID: 1, Amount: decimal.NewFromFloat(100.50), Status: "success"},
		{UserID: 1, Amount: decimal.NewFromFloat(200.00), Status: "pending"},
		{UserID: 1, Amount: decimal.NewFromFloat(150.75), Status: "success"},
		{UserID: 2, Amount: decimal.NewFromFloat(300.00), Status: "success"},
	}

	for _, tx := range transactions {
		suite.repo.Create(tx)
	}

	// Test with all filters
	filters := models.TransactionFilters{
		UserID: 1,
		Status: "success",
		Limit:  1,
		Offset: 0,
	}

	result, err := suite.repo.GetAll(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1) // Limited to 1
	assert.Equal(suite.T(), uint(1), result[0].UserID)
	assert.Equal(suite.T(), "success", result[0].Status)
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
