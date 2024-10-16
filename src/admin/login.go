package admin

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/mviner000/eyygo/src/admin/models"
	"github.com/mviner000/eyygo/src/auth"
	"github.com/mviner000/eyygo/src/config"
	"github.com/mviner000/eyygo/src/http"
)

func LoginForm(c *fiber.Ctx) error {
	log.Println("LoginForm function called")
	errorMessage := c.Query("error")
	// Check if it's an HTMX request
	if c.Get("HX-Request") == "true" {
		log.Println("HTMX request detected")
		// If it's an HTMX request, just return the form content
		return http.HttpResponseHTMX(fiber.Map{
			"Error": errorMessage,
		}, "src/admin/templates/login_form.html").Render(c)
	}
	log.Println("Rendering full login page")
	// Render the full page with layout
	return http.HttpResponseHTMX(fiber.Map{
		"Error":     errorMessage,
		"MetaTitle": "Login | " + SiteName,
	}, "src/admin/templates/login.html", "src/admin/templates/layout.html").Render(c)
}

func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	authUser, err := auth.GetUserByUsername(username)
	if err != nil {
		log.Printf("Login attempt failed for username '%s': user not found", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	// Use the password verification method
	match, err := auth.VerifyPassword(authUser.Password, password)
	if err != nil || !match {
		log.Printf("Login attempt failed for user '%s': invalid password", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	// Use the JWT token generator
	tokenGenerator := auth.NewPasswordResetTokenGenerator()
	token, err := tokenGenerator.MakeToken(authUser)
	if err != nil {
		log.Printf("Failed to generate JWT token for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to generate authentication token", nil).Render(c)
	}
	log.Printf("Generated JWT token for user '%s'", username)

	// Update last_login in the database
	if err := auth.UpdateLastLogin(authUser.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", username, err)
	}

	// Create database session
	sessionID := generateSessionID()
	expireTime := time.Now().Add(24 * time.Hour)
	session := models.Session{
		SessionKey: sessionID,
		UserID:     uint(authUser.ID),
		AuthToken:  token,
		ExpireDate: expireTime,
	}
	if err := config.GetDB().Create(&session).Error; err != nil {
		log.Printf("Failed to create session for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to create session", nil).Render(c)
	}
	log.Printf("Created session for user '%s': sessionID=%s", username, sessionID)

	// Calculate maxAge using time.Until
	maxAge := int(time.Until(expireTime).Seconds())
	// Set the browser client session cookie using the utility function
	auth.SetSessionCookie(c, sessionID, expireTime, maxAge)

	log.Printf("Login successful for user '%s'. Session ID: %s", username, sessionID)
	return c.SendString(http.WindowReload("/admin/dashboard"))
}

// generateSessionID generates a new session ID
func generateSessionID() string {
	return uuid.New().String()
}
