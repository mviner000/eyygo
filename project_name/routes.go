package project_name

import (
	"github.com/gofiber/fiber/v2"
)

// AppName implements the App interface
type AppName struct{}

// SetupRoutes sets up the routes for the project_name app
func (a *AppName) SetupRoutes(app *fiber.App) {
	// Set up routes for the project_name app
	app.Get("/project_name", func(c *fiber.Ctx) error {
		return c.SendString("Hello from project_name!")
	})
}
