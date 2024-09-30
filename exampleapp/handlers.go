package exampleapp

import (
	"github.com/gofiber/fiber/v2"
)

// HelloHandler handles the hello route
func HelloHandler(c *fiber.Ctx) error {
	return c.SendString("Hello from exampleapp!")
}
