package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecoveryMiddlewarePanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
