package handlers

import (
	"strconv"

	"interview/internal/models"
	"interview/internal/services"
	"interview/pkg/utils"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TransactionHandler handles transaction HTTP requests
type TransactionHandler struct {
	service   services.TransactionService
	validator *validator.Validate
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	validator := validator.New()

	// Register custom validation for decimal.Decimal
	validator.RegisterValidation("decimal_positive", validateDecimalPositive)

	return &TransactionHandler{
		service:   service,
		validator: validator,
	}
}

// validateDecimalPositive validates that a decimal.Decimal value is positive
func validateDecimalPositive(fl validator.FieldLevel) bool {
	amount := fl.Field().Interface().(decimal.Decimal)
	return amount.GreaterThan(decimal.Zero)
}

// CreateTransaction handles POST /api/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.BadRequestResponse(c, "Validation failed: "+err.Error())
		return
	}

	transaction, err := h.service.CreateTransaction(req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, transaction, "Transaction created successfully")
}

// GetTransactions handles GET /api/transactions
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	var filters models.TransactionFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.BadRequestResponse(c, "Invalid query parameters")
		return
	}

	transactions, err := h.service.GetTransactions(filters)
	if err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, transactions, "Transactions retrieved successfully")
}

// GetTransaction handles GET /api/transactions/:id
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid transaction ID")
		return
	}

	transaction, err := h.service.GetTransaction(uint(id))
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.NotFoundResponse(c, "Transaction not found")
			return
		}
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, transaction, "Transaction retrieved successfully")
}

// UpdateTransaction handles PUT /api/transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid transaction ID")
		return
	}

	var req models.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.BadRequestResponse(c, "Validation failed: "+err.Error())
		return
	}

	err = h.service.UpdateTransactionStatus(uint(id), req.Status)
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.NotFoundResponse(c, "Transaction not found")
			return
		}
		if err.Error() == "invalid status" {
			utils.BadRequestResponse(c, "Invalid status")
			return
		}
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Transaction updated successfully")
}

// DeleteTransaction handles DELETE /api/transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid transaction ID")
		return
	}

	err = h.service.DeleteTransaction(uint(id))
	if err != nil {
		if err.Error() == "transaction not found" {
			utils.NotFoundResponse(c, "Transaction not found")
			return
		}
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "Transaction deleted successfully")
}
