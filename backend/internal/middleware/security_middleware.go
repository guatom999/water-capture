package middleware

import (
	"github.com/labstack/echo/v4"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Prevent MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Content Security Policy
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")

			// Prevent XSS attacks
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")

			return next(c)
		}
	}
}
