package admin

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the admin routes
func SetupRoutes(app *fiber.App) {
	// Define a group for admin routes
	adminGroup := app.Group("/admin")

	// Define the login routes
	adminGroup.Get("/login", LoginForm)
	adminGroup.Post("/login", Login)

	// Define the dashboard route
	adminGroup.Get("/dashboard", Dashboard)

	// Define user management routes
	adminGroup.Get("/users", UserIndex)
	adminGroup.Get("/users/create", UserCreate)
	adminGroup.Post("/users", UserStore)
}
