// main.go
package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/mviner000/eyygo/config"
	"github.com/mviner000/eyygo/handlers"
	"github.com/mviner000/eyygo/logger"
	"github.com/mviner000/eyygo/middleware"
	"github.com/mviner000/eyygo/models"
	"github.com/mviner000/eyygo/routes"
	"github.com/mviner000/eyygo/settings"
	"github.com/mviner000/eyygo/views"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

func main() {
	// Initialize logger
	appLogger := logger.NewLogger()

	// 1. First load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize database connection and store in global variable
	DB, err = settings.NewDBConnection(cfg)
	if err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Test database connection
	sqlDB, err := DB.DB()
	if err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize JWT secret
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-secret-key") // Default secret - change in production
	}

	// 4. Initialize handlers (after DB is initialized)
	authHandler := handlers.NewAuthHandler(DB, jwtSecret)
	adminHandler := handlers.NewAdminHandler(DB)
	viewHandler := views.NewViewHandler(DB)

	// 5. Auto-migrate the database
	if err := DB.AutoMigrate(&models.User{}, &models.Note{}); err != nil {
		appLogger.ErrorLogger.Printf("Failed to auto-migrate: %v", err)
	}

	// Initialize Fiber app
	// Initialize Fiber app with template engine
	app := fiber.New(fiber.Config{
		Views: html.New("./templates", ".html"), // Set template engine
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			appLogger.ErrorLogger.Printf("Error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		},
	})

	// Global Middleware
	app.Use(middleware.XFrameOptions())    // Clickjacking protection
	app.Use(middleware.ConfigureCORS(cfg)) // CORS with config
	app.Use(middleware.RateLimit())        // Rate limiting
	app.Use(middleware.SecurityHeaders())  // Additional security headers
	app.Use(recover.New())                 // Recover from panics
	app.Use(logger.RequestLogger())        // Request logging
	app.Use(logger.ErrorLogger())          // Error logging

	// Setup routes with viewHandler
	routes.SetupRoutes(app, authHandler, adminHandler, viewHandler, jwtSecret)

	// Print server status (Django-style)
	appLogger.PrintServerStatus(cfg.ServerHost, cfg.ServerPort)

	// Start server
	if err := app.Listen(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
