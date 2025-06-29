package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
