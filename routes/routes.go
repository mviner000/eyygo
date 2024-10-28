// routes/routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyygo/handlers"
	"github.com/mviner000/eyygo/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, adminHandler *handlers.AdminHandler, jwtSecret []byte) {
	// Public routes
	setupPublicRoutes(app, authHandler)

	// Protected API routes
	api := app.Group("/api")
	api.Use(middleware.Protected(jwtSecret))
	setupAPIRoutes(api, authHandler)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(handlers.AdminMiddleware)
	setupAdminRoutes(admin, adminHandler)
}

// setupPublicRoutes configures public routes
func setupPublicRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	app.Post("/auth/login", authHandler.Login)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     true,
		})
	})

	app.Get("/db-test", func(c *fiber.Ctx) error {
		return c.SendString("Database connection successful!")
	})
}

// setupAPIRoutes configures protected API routes
func setupAPIRoutes(api fiber.Router, authHandler *handlers.AuthHandler) {
	api.Get("/auth/validate", authHandler.ValidateToken)
}

// setupAdminRoutes configures admin routes
func setupAdminRoutes(admin fiber.Router, adminHandler *handlers.AdminHandler) {
	admin.Get("/users", adminHandler.ListUsers)
	admin.Post("/users", adminHandler.CreateUser)
	admin.Put("/users/:id", adminHandler.UpdateUser)
	admin.Delete("/users/:id", handlers.SuperUserMiddleware, adminHandler.DeleteUser)
}
