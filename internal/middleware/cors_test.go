package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"interview/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
}

func TestCORSMiddlewareOptionsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}
