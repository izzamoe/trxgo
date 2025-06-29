package utils

import (
	"net/http"

	"interview/internal/models"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := models.APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusOK, response)
}

// CreatedResponse sends a created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	response := models.APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusCreated, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	response := models.APIResponse{
		Success: false,
		Error:   message,
	}
	c.JSON(statusCode, response)
}

// BadRequestResponse sends a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}
