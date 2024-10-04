package decorators

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/types"
	"gorm.io/gorm"
)

// RequireAuthentication ensures the user is authenticated
func RequireAuthentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Implement your authentication check here
		// This is just a placeholder implementation
		if c.Locals("user") == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}
		return c.Next()
	}
}

// RequireSuperuser ensures the user is a superuser
func RequireSuperuser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*types.User)
		if !ok || !user.IsSuperuser {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Superuser access required",
			})
		}
		return c.Next()
	}
}

// DatabaseTransaction wraps the handler in a database transaction
func DatabaseTransaction(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return db.Transaction(func(tx *gorm.DB) error {
			c.Locals("tx", tx)
			return c.Next()
		})
	}
}

// Logger logs incoming requests
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.DebugLog("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	}
}

// Throttle limits the number of requests from a single IP
func Throttle(limit int, duration int) fiber.Handler {
	// Implement rate limiting logic here
	// This is just a placeholder
	return func(c *fiber.Ctx) error {
		// Implement rate limiting logic
		return c.Next()
	}
}
