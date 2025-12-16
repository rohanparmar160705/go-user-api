/*
Package middleware provides HTTP middleware functions.
RequestID middleware generates a unique UUID for each incoming request.
This ID is added to the response headers (X-Request-ID) and the context,
allowing for request tracing across logs.
*/
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a new UUID for the request
		requestID := uuid.New().String()
		
		// Set request ID in context (for logging)
		c.Locals("requestID", requestID)
		
		// Add request ID to response header
		c.Set("X-Request-ID", requestID)
		
		// Continue to next handler
		return c.Next()
	}
}
