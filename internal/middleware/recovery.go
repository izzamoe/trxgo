package middleware

import (
	"net/http"

	"interview/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithFields(logrus.Fields{
			"error":  recovered,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Error("Panic recovered")

		utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	})
}
