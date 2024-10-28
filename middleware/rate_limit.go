package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimit creates a rate limiting middleware
func RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // max requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // use IP address as key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}
