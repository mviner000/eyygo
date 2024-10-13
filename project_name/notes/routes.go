// notes/routes.go
package notes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupNoteRoutes sets up all the note-related routes under the provided group
func SetupNoteRoutes(app fiber.Router) {
	// Public Note Routes

	// List all notes
	app.Get("/", listNotes)

	// Retrieves a note by its ID
	app.Get("/:id", getNoteByID)

	// Create Note Route
	app.Post("/", createNote)

	// Delete Note Route
	app.Delete("/:id", deleteNote)
}
