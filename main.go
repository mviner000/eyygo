package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/mviner000/eyygo/config"
	"github.com/mviner000/eyygo/handlers"
	"github.com/mviner000/eyygo/logger"
	"github.com/mviner000/eyygo/models"
	"github.com/mviner000/eyygo/settings"
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

	// 5. Auto-migrate the database
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		appLogger.ErrorLogger.Printf("Failed to auto-migrate: %v", err)
	}

	// 6. Create default superuser if not exists
	var superUser models.User
	if err := DB.Where("username = ?", "admin").First(&superUser).Error; err != nil {
		if err := models.CreateSuperUser(DB, "admin", "admin@example.com", "adminpassword"); err != nil {
			appLogger.ErrorLogger.Printf("Failed to create superuser: %v", err)
		}
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			appLogger.ErrorLogger.Printf("Error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		},
	})

	// Public routes
	app.Post("/auth/login", authHandler.Login)

	// Middleware
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.RequestLogger())
	app.Use(logger.ErrorLogger())

	// JWT Middleware - apply only to protected routes
	app.Use("/admin", jwtware.New(jwtware.Config{
		SigningKey: jwtSecret,
	}))

	// Setup routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Validate token endpoint
	app.Get("/auth/validate", jwtware.New(jwtware.Config{
		SigningKey: jwtSecret,
	}), authHandler.ValidateToken)

	// Add test route for database connection
	app.Get("/db-test", func(c *fiber.Ctx) error {
		result := DB.Raw("SELECT 1")
		if result.Error != nil {
			return c.Status(500).SendString("Database connection failed")
		}
		return c.SendString("Database connection successful!")
	})

	// Admin routes
	admin := app.Group("/admin")
	admin.Use(handlers.AdminMiddleware)

	// User management routes
	admin.Get("/users", adminHandler.ListUsers)
	admin.Post("/users", adminHandler.CreateUser)
	admin.Put("/users/:id", adminHandler.UpdateUser)
	admin.Delete("/users/:id", handlers.SuperUserMiddleware, adminHandler.DeleteUser)

	// Print server status (Django-style)
	appLogger.PrintServerStatus(cfg.ServerHost, cfg.ServerPort)

	// Start server
	if err := app.Listen(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		appLogger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
