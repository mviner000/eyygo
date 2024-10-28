package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyygo/handlers"
	"github.com/mviner000/eyygo/middleware"
	"github.com/mviner000/eyygo/views"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, adminHandler *handlers.AdminHandler, viewHandler *views.ViewHandler, jwtSecret []byte) {
	// Public routes
	setupPublicRoutes(app, authHandler, viewHandler)

	// Protected routes
	protected := app.Group("/")
	protected.Use(middleware.Protected(jwtSecret))
	setupProtectedRoutes(protected, viewHandler)

	// Protected API routes
	api := app.Group("/api")
	api.Use(middleware.Protected(jwtSecret))
	setupAPIRoutes(api, authHandler)

	// Admin routes
	admin := app.Group("/admin")
	setupAdminRoutes(admin, authHandler, adminHandler)

	// Protected admin API routes
	adminAPI := api.Group("/admin")
	setupAdminAPIRoutes(adminAPI, adminHandler)
}

// setupPublicRoutes configures public routes
func setupPublicRoutes(app fiber.Router, authHandler *handlers.AuthHandler, viewHandler *views.ViewHandler) {
	// Static pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"db":     true,
		})
	})

	// Authentication routes
	app.Get("/login", viewHandler.LoginPage)
	app.Post("/validate-username", viewHandler.ValidateUsername)
	app.Post("/validate-password", viewHandler.ValidatePassword)
	app.Post("/validate-login", viewHandler.ValidateLogin)
	app.Get("/logout", viewHandler.LogoutPage)
}

// setupProtectedRoutes configures routes that require authentication
func setupProtectedRoutes(router fiber.Router, viewHandler *views.ViewHandler) {
	router.Get("/dashboard", viewHandler.Dashboard)
	router.Get("/users/list", viewHandler.UsersList)
	router.Get("/notes/list", viewHandler.NotesList)
}

// setupAPIRoutes configures protected API routes
func setupAPIRoutes(api fiber.Router, authHandler *handlers.AuthHandler) {
	api.Get("/auth/validate", authHandler.ValidateToken)
}

// setupAdminRoutes configures admin panel routes
func setupAdminRoutes(admin fiber.Router, authHandler *handlers.AuthHandler, adminHandler *handlers.AdminHandler) {
	admin.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("admin/login", fiber.Map{
			"title": "Admin Login",
		})
	})
	admin.Post("/login", authHandler.Login)
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

// NewRoutes initializes all routes
func NewRoutes(db *gorm.DB, jwtSecret []byte) (*handlers.AuthHandler, *handlers.AdminHandler, *views.ViewHandler) {
	authHandler := handlers.NewAuthHandler(db, jwtSecret)
	adminHandler := handlers.NewAdminHandler(db)
	viewHandler := views.NewViewHandler(db)

	return authHandler, adminHandler, viewHandler
}
