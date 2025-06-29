package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/pkg/utils"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"key": "value"}
	utils.SuccessResponse(c, testData, "Success message")

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}
}

func TestCreatedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"id": "123"}
	utils.CreatedResponse(c, testData, "Created successfully")

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.ErrorResponse(c, http.StatusBadRequest, "Bad request error")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBadRequestResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.BadRequestResponse(c, "Invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestNotFoundResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.NotFoundResponse(c, "Resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestInternalServerErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	utils.InternalServerErrorResponse(c, "Server error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
