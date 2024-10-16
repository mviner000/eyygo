package monitor

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all the monitoring routes
func SetupRoutes(app *fiber.App) {
	app.Get("/status", HandleStatus)
	app.Get("/status/server-info", HandleStatus)
	app.Get("/status/cpu", HandleStatus)
	app.Get("/status/ram", HandleStatus)
	app.Get("/status/storage", HandleStatus)
	app.Get("/status/old", HandleStatus)
}
