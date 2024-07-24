package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// LogrusLogger is a middleware function that logs HTTP requests using Logrus
func LogrusLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log details
		logger.WithFields(logrus.Fields{
			"status_code":   statusCode,
			"latency_time":  latency,
			"client_ip":     c.ClientIP(),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"user_agent":    c.Request.UserAgent(),
			"error_message": c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}).Debug("HTTP request log")
	}
}
