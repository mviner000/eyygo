package exampleapp

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures the routes for the exampleapp app
func SetupRoutes(app *fiber.App) {
	exampleappGroup := app.Group("/exampleapp")
	exampleappGroup.Get("/hello", HelloHandler)
}
