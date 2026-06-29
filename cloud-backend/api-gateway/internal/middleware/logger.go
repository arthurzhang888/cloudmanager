package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)

		// Get user info if available
		userID, _ := c.Get("user_id")
		if userID == nil {
			userID = "anonymous"
		}

		// Log format: timestamp | request_id | method | path | status | duration | user_id
		gin.DefaultWriter.Write([]byte(
			time.Now().Format("2006-01-02 15:04:05") + " | " +
				requestID + " | " +
				c.Request.Method + " | " +
				c.Request.URL.Path + " | " +
				duration.String() + " | " +
				userID.(string) + "\n",
		))
	}
}

// Cors middleware handles CORS headers
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler handles panics and returns proper error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				if requestID == nil {
					requestID = "unknown"
				}

				c.JSON(500, gin.H{
					"error":      "internal server error",
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}
