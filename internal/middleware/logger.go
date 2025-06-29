package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware provides request logging
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"error":       param.ErrorMessage,
			"body_size":   param.BodySize,
			"user_agent":  param.Request.UserAgent(),
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")

		return ""
	})
}
