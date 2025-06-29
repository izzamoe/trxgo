package services

import (
	"fmt"

	"interview/internal/models"
	"interview/internal/repositories"
)

// DashboardService interface defines dashboard service methods
type DashboardService interface {
	GetSummary() (*models.DashboardSummary, error)
}

// dashboardService implements DashboardService interface
type dashboardService struct {
	repo repositories.TransactionRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(repo repositories.TransactionRepository) DashboardService {
	return &dashboardService{repo: repo}
}

// GetSummary gets dashboard summary
func (s *dashboardService) GetSummary() (*models.DashboardSummary, error) {
	summary := &models.DashboardSummary{}

	// Get today's successful transactions
	todayCount, todayAmount, err := s.repo.GetTodaySuccessful()
	if err != nil {
		return nil, fmt.Errorf("failed to get today's successful transactions: %v", err)
	}
	summary.TodaySuccessfulTransactions = todayCount
	summary.TodaySuccessfulAmount = todayAmount

	// Get average transactions per user
	avgPerUser, err := s.repo.GetAveragePerUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get average transactions per user: %v", err)
	}
	summary.AverageTransactionPerUser = avgPerUser

	// Get latest transactions
	latestTransactions, err := s.repo.GetLatest(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest transactions: %v", err)
	}
	summary.LatestTransactions = latestTransactions

	// Get status counts
	statusCounts, err := s.repo.GetStatusCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %v", err)
	}
	summary.StatusCounts = statusCounts

	return summary, nil
}
