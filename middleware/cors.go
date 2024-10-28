package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mviner000/eyygo/config"
)

// ConfigureCORS returns CORS middleware with custom configuration
func ConfigureCORS(cfg *config.Config) fiber.Handler {
	// Determine allowed origins based on environment
	allowedOrigins := "http://" + cfg.ServerHost + ":" + cfg.ServerPort

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 hours
	})
}
