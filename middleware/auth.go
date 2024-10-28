package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// Protected creates a middleware that verifies JWT tokens
func Protected(jwtSecret []byte) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtSecret,
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid or expired token",
	})
}
