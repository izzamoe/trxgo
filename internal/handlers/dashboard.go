package handlers

import (
	"interview/internal/services"
	"interview/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard HTTP requests
type DashboardHandler struct {
	service services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// GetSummary handles GET /api/dashboard/summary
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, summary, "Dashboard summary retrieved successfully")
}
