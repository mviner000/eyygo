package admin

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mviner000/eyymi/eyygo/auth"
	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/http"

	models "github.com/mviner000/eyymi/eyygo/admin/models"
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
		}, "eyygo/admin/templates/login_form.html").Render(c)
	}

	log.Println("Rendering full login page")
	// Render the full page with layout
	return http.HttpResponseHTMX(fiber.Map{
		"Error":     errorMessage,
		"MetaTitle": "Login | " + SiteName,
	}, "eyygo/admin/templates/login.html", "eyygo/admin/templates/layout.html").Render(c)
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

	match, err := config.CheckPasswordHash(password, authUser.Password)
	if err != nil || !match {
		log.Printf("Login attempt failed for user '%s': invalid password", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	token, err := tokenGenerator.MakeToken(authUser)

	if err != nil {
		log.Printf("Failed to generate authentication token for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to generate authentication token", nil).Render(c)
	}
	log.Printf("Generated token for user '%s': %s", username, token)

	// Update last_login in the database
	if err := auth.UpdateLastLogin(authUser.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", username, err)
	}

	// Create session
	sessionID := generateSessionID()
	expireTime := time.Now().Add(24 * time.Hour)

	session := models.Session{
		SessionKey: sessionID,
		UserID:     uint(authUser.ID), // Convert int to uint
		AuthToken:  token,
		ExpireDate: expireTime,
	}

	if err := config.GetDB().Create(&session).Error; err != nil {
		log.Printf("Failed to create session for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to create session", nil).Render(c)
	}
	log.Printf("Created session for user '%s': sessionID=%s, token=%s", username, sessionID, token)

	// Calculate maxAge using time.Until
	maxAge := int(time.Until(expireTime).Seconds())

	// Set the session cookie using the new utility function
	auth.SetSessionCookie(c, sessionID, expireTime, maxAge)

	log.Printf("Login successful for user '%s'. Session ID: %s, Auth Token: %s", username, sessionID, token)

	return c.SendString(http.WindowReload("/admin/dashboard"))
}

// generateSessionID generates a new session ID
func generateSessionID() string {
	return uuid.New().String()
}
