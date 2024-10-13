// notes/routes.go
package notes

import (
	"github.com/gofiber/fiber/v2"
	auth "github.com/mviner000/eyymi/eyygo/auth"
)

// SetupNoteRoutes sets up all the note-related routes under the provided group
func SetupNoteRoutes(app fiber.Router) {
	// Apply JWT middleware to all note routes

	// Initialize the secret key for JWT verification
	auth.InitJWTSecret()
	app.Use(auth.JWTMiddleware())

	// CRUD operations for notes
	app.Post("/", createNote)      // Create a new note
	app.Get("/", listNotes)        // List all notes
	app.Get("/:id", getNoteByID)   // Retrieve a note by its ID
	app.Delete("/:id", deleteNote) // Delete a note
}
