// notes/views.go
package notes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/src/config"
	"gorm.io/gorm"
)

// Create Note Handler
func createNote(c *fiber.Ctx) error {
	note := new(Note)

	// Parse the JSON request body into the note struct
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Set the created and updated timestamps
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	// Save the note to the database
	db := config.GetDB() // Assuming you have a GetDB() function to obtain the GORM DB instance
	if err := db.Create(note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create note",
		})
	}

	// Return a successful response with the created note
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Note created successfully",
		"note":    note,
	})
}

// List Notes Handler
func listNotes(c *fiber.Ctx) error {
	var notes []Note

	// Retrieve the database connection

	db := getDB()
	// Fetch all notes from the database
	if err := db.Find(&notes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve notes",
		})
	}

	// Return a successful response with the list of notes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notes": notes,
	})
}

// Get Note by ID Handler
func getNoteByID(c *fiber.Ctx) error {
	id := c.Params("id") // Extract the note ID from the URL
	var note Note

	// Retrieve the database connection
	db := getDB()

	// Find the note by ID
	if err := db.First(&note, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve note",
		})
	}

	// Return the retrieved note
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"note": note,
	})
}

// Delete Note Handler
func deleteNote(c *fiber.Ctx) error {
	id := c.Params("id") // Extract the note ID from the URL

	// Retrieve the database connection
	db := getDB()

	// Attempt to delete the note
	result := db.Delete(&Note{}, id)

	// Check for any errors during the delete operation
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Note not found",
		})
	}

	// Return a success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}

// Utility function to get the database connection
func getDB() *gorm.DB {
	// Assuming you have a config package that provides the database connection
	return config.GetDB()
}
