/*
Package middleware provides HTTP middleware functions.
RequestDuration middleware measures the time taken to process each request.
It logs the duration, HTTP status, method, and path using the structured logger context.
Must be used after RequestID middleware to include the request ID in logs.
*/
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rohanparmar/go-user-api/internal/logger"
	"go.uber.org/zap"
)

// RequestDuration middleware logs the duration of each request
func RequestDuration() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Record start time
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Get request ID from context
		requestID, _ := c.Locals("requestID").(string)
		
		// Log request details with duration
		logger.Log.Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)
		
		return err
	}
}
