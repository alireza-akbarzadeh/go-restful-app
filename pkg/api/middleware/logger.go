package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger creates a logging middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)

		// Get status
		statusCode := c.Writer.Status()

		// Log format
		log.Printf("[%s] %s %s %d %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}
