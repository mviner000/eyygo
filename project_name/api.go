// project_name/api.go
package project_name

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/project_name/notes"
)

// SetupAPIRoutes sets up all the API routes under the /api prefix
func SetupAPIRoutes(app *fiber.App) {

	// Initialize the secret key for JWT verification
	notes.InitJWTSecret()

	// Group all API routes under /api
	apiGroup := app.Group("/api")

	// Public API routes
	apiGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to project_name API!")
	})

	// Token generation endpoint
	apiGroup.Post("/token/pair", generateTokenPairHandler)

	// Group notes-related routes under /api/notes
	noteGroup := apiGroup.Group("/notes")

	// Call the function to set up note routes
	notes.SetupNoteRoutes(noteGroup) // Pass the noteGroup to the SetupNoteRoutes function
}

func generateTokenPairHandler(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if username is "test_username"
	if body.Username != "test_username" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username",
		})
	}

	// Generate token pair
	accessToken, refreshToken, err := GenerateTokenPair(body.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate tokens",
		})
	}

	// Return token pair
	return c.JSON(fiber.Map{
		"token_type":    "Bearer",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
