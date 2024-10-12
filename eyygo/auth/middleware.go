package auth

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/config"
	"gorm.io/gorm"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(c *fiber.Ctx) error {
	log.Printf("AuthMiddleware: Starting authentication check for request: %s %s", c.Method(), c.Path())

	// Check if the current path is "/admin/login"
	if c.Path() == "/admin/login" {
		// Set maxAge to 0 to expire the session cookie
		SetSessionCookie(c, "", time.Now(), 0)
		DeleteSessionCookie(c)
		log.Println("AuthMiddleware: Session cookie expired and deleted for login path.")
		return c.Next()
	}

	// Retrieve session from the database
	userIDStr, authToken, err := getSessionFromDB(c)
	if err != nil {
		log.Printf("AuthMiddleware: Error retrieving session: %v", err)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: Session retrieved for user ID: %s, Token: %s", userIDStr, authToken)

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("AuthMiddleware: Invalid user ID: %s", userIDStr)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: User ID converted to int: %d", userID)

	// Get the user from the database
	user, err := GetUserByID(uint(userID))
	if err != nil {
		log.Printf("AuthMiddleware: Error retrieving user from database: %v", err)
		return c.Redirect("/admin/login")
	}
	if user == nil {
		log.Printf("AuthMiddleware: User not found for ID: %d", userID)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: User retrieved from database: %s", user.Username)

	// Check if the token is valid for the user
	tokenGenerator := NewPasswordResetTokenGenerator()
	if !tokenGenerator.CheckToken(user, authToken) {
		log.Printf("AuthMiddleware: Invalid token for user %s. Token: %s", user.Username, authToken)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: Token valid for user %s", user.Username)

	// Store user information in the context for later use
	c.Locals("user", user)
	log.Printf("AuthMiddleware: User %s stored in context", user.Username)

	// All checks passed, proceed to the next handler
	log.Println("AuthMiddleware: Authentication successful, proceeding to next handler")
	return c.Next()
}

func getSessionFromDB(c *fiber.Ctx) (string, string, error) {
	log.Println("getSessionFromDB: Starting session retrieval")

	// Get session ID from cookie
	sessionID := c.Cookies(SessionCookieName)
	if sessionID == "" {
		log.Println("getSessionFromDB: Session ID not found in cookie")
		return "", "", fmt.Errorf("session ID not found in cookie")
	}
	log.Printf("getSessionFromDB: Session ID found: %s", sessionID)

	db := config.GetDB()
	var session models.Session
	result := db.Where("session_key = ?", sessionID).First(&session)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("getSessionFromDB: Session not found")
			return "", "", fmt.Errorf("session not found")
		}
		log.Printf("getSessionFromDB: Error querying database: %v", result.Error)
		return "", "", result.Error
	}

	// Check if the session is expired
	if session.ExpireDate.Before(time.Now()) {
		log.Println("getSessionFromDB: Session found but expired")
		return "", "", fmt.Errorf("session expired")
	}

	log.Printf("getSessionFromDB: Session retrieved for user ID: %d, Token: %s", session.UserID, session.AuthToken)
	return fmt.Sprintf("%d", session.UserID), session.AuthToken, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(userID uint) (*models.AuthUser, error) {
	log.Printf("GetUserByID: Retrieving user with ID %d", userID)

	db := config.GetDB()
	var user models.AuthUser
	result := db.First(&user, userID)
	if result.Error != nil {
		log.Printf("GetUserByID: Error retrieving user by ID %d from database: %v", userID, result.Error)
		return nil, result.Error
	}

	log.Printf("GetUserByID: User ID %d (%s) retrieved successfully from database", userID, user.Username)
	return &user, nil
}
