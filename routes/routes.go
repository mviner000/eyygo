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
	admin := app.Group("/admin")
	setupAdminLoginRoutes(admin, authHandler)

	// Protected admin API routes
	adminAPI := api.Group("/admin")
	setupAdminAPIRoutes(adminAPI, adminHandler)
}

// setupPublicRoutes configures public routes
func setupPublicRoutes(app fiber.Router, authHandler *handlers.AuthHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     true,
		})
	})
}

// setupAdminLoginRoutes configures admin login routes
func setupAdminLoginRoutes(admin fiber.Router, authHandler *handlers.AuthHandler) {
	admin.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("admin/login", fiber.Map{
			"title": "Admin Login",
		})
	})
	admin.Post("/login", authHandler.Login)
}

// setupAPIRoutes configures protected API routes
func setupAPIRoutes(api fiber.Router, authHandler *handlers.AuthHandler) {
	api.Get("/auth/validate", authHandler.ValidateToken)
}

// setupAdminAPIRoutes configures protected admin API routes
func setupAdminAPIRoutes(admin fiber.Router, adminHandler *handlers.AdminHandler) {
	admin.Get("/models", adminHandler.ListModels)
	admin.Get("/models/:model", adminHandler.ListModelEntries)
	admin.Get("/models/:model/:id", adminHandler.GetModelEntry)
	admin.Post("/models/:model", adminHandler.CreateModelEntry)
	admin.Put("/models/:model/:id", adminHandler.UpdateModelEntry)
	admin.Delete("/models/:model/:id", adminHandler.DeleteModelEntry)
}
