// notes/routes.go
package notes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupNoteRoutes sets up all the note-related routes under the provided group
func SetupNoteRoutes(app fiber.Router) {
	// Apply JWT middleware to all note routes
	app.Use(JWTMiddleware())

	// CRUD operations for notes
	app.Post("/", createNote)      // Create a new note
	app.Get("/", listNotes)        // List all notes
	app.Get("/:id", getNoteByID)   // Retrieve a note by its ID
	app.Delete("/:id", deleteNote) // Delete a note
}
