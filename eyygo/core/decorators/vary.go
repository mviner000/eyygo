package decorators

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// VaryOnHeaders adds the Vary header with the specified headers
func VaryOnHeaders(headers ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Vary(headers...)
		return c.Next()
	}
}

// CacheControlMaxAge sets the Cache-Control header with a max-age directive
func CacheControlMaxAge(maxAge int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
		return c.Next()
	}
}
