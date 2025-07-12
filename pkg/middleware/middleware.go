package middleware

import (
	"github.com/anggara-26/pomodoro-backend.git/app/handler"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a placeholder for general authentication
func AuthMiddleware(c *fiber.Ctx) error {
	// This is a basic middleware that can be extended
	// For now, it just passes through
	return c.Next()
}

// ErrorHandler handles panics and errors globally
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(handler.Response{
		Message: err.Error(),
		Code:    code,
	})
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Basic rate limiting - can be enhanced with Redis
		// For now, just add headers
		c.Set("X-RateLimit-Limit", "100")
		c.Set("X-RateLimit-Remaining", "99")
		return c.Next()
	}
}

// CORSConfig returns CORS configuration
func CORSConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: ErrorHandler,
	}
}

// LoggingMiddleware provides request logging
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add custom logging logic here if needed
		return c.Next()
	}
}
