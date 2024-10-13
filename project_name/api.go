// project_name/api.go
package project_name

import (
	"github.com/gofiber/fiber/v2"
	notes "github.com/mviner000/eyymi/project_name/notes" // Import the notes package
)

// SetupAPIRoutes sets up all the API routes under the /api prefix
func SetupAPIRoutes(app *fiber.App) {
	// Group all API routes under /api
	apiGroup := app.Group("/api")

	// Public API routes
	apiGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to project_name API!")
	})

	// Group notes-related routes under /api/notes
	noteGroup := apiGroup.Group("/notes")

	// Call the function to set up note routes
	notes.SetupNoteRoutes(noteGroup) // Pass the noteGroup to the SetupNoteRoutes function
}
