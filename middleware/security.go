// middleware/security.go
package middleware

import "github.com/gofiber/fiber/v2"

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-Frame-Options for clickjacking protection
		c.Set("X-Frame-Options", "DENY")

		// XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// Content type options
		c.Set("X-Content-Type-Options", "nosniff")

		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		c.Set("Content-Security-Policy", `
            default-src 'self';
            script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com https://cdn.tailwindcss.com;
            style-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com;
            img-src 'self' data: https:;
            font-src 'self' data: https:;
            connect-src 'self' https://unpkg.com https://cdn.tailwindcss.com;
        `)

		return c.Next()
	}
}

func XFrameOptions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		return c.Next()
	}
}
