// project_name/routes.go
package project_name

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/src/auth"
)

// AppName implements the App interface
type AppName struct{}

// SetupRoutes sets up the routes for the project_name app
func (a *AppName) SetupRoutes(app *fiber.App) {
	log.Println("Admin: Starting to set up routes")

	// Set up non-API routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the project_name homepage!")
	})

	// Group admin-related routes under /api/admin
	adminGroup := app.Group("/api/admin")

	// Protected Admin Routes
	adminGroup.Get("/", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("This is an admin route!")
	})

	adminGroup.Get("/dashboard", auth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("This is the admin dashboard!")
	})

	// Call the function to set up API routes
	SetupAPIRoutes(app)
}
