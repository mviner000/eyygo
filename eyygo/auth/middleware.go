package auth

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(c *fiber.Ctx) error {
	log.Println("AuthMiddleware: Starting authentication check")

	// Retrieve session from the database
	userIDStr, authToken, err := getSessionFromDB(c)
	if err != nil {
		log.Printf("AuthMiddleware: Error retrieving session: %v", err)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: Session retrieved for user ID: %s", userIDStr)

	// Convert userID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("AuthMiddleware: Invalid user ID: %s", userIDStr)
		return c.Redirect("/admin/login")
	}
	log.Printf("AuthMiddleware: User ID converted to int: %d", userID)

	// Get the user from the database
	user, err := GetUserByID(userID)
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
		log.Printf("AuthMiddleware: Invalid token for user %s, redirecting to login", user.Username)
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

	var userID string
	var authToken string

	// Get session ID from cookie
	sessionID := c.Cookies("hey_sesion")
	if sessionID == "" {
		log.Println("getSessionFromDB: Session ID not found in cookie")
		return "", "", fmt.Errorf("session ID not found in cookie")
	}
	log.Printf("getSessionFromDB: Session ID found: %s", sessionID)

	// Check if session exists without considering expiration
	query := `SELECT user_id, auth_token, expire_date FROM eyygo_session WHERE session_key = ?`
	var expireDate time.Time
	err := db.QueryRow(query, sessionID).Scan(&userID, &authToken, &expireDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("getSessionFromDB: Session not found")
			return "", "", fmt.Errorf("session not found")
		}
		log.Printf("getSessionFromDB: Error querying database: %v", err)
		return "", "", err
	}

	// Check if the session is expired
	if expireDate.Before(time.Now()) {
		log.Println("getSessionFromDB: Session found but expired")
		return "", "", fmt.Errorf("session expired")
	}

	log.Printf("getSessionFromDB: Session retrieved for user ID: %s", userID)
	return userID, authToken, nil
}

// GetUserByID retrieves a user by ID from the database
func GetUserByID(userID int) (*User, error) {
	log.Printf("GetUserByID: Retrieving user with ID %d", userID)

	var user User
	query := "SELECT id, username, email, password, last_login FROM auth_user WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.LastLogin)
	if err != nil {
		log.Printf("GetUserByID: Error retrieving user by ID %d from database: %v", userID, err)
		return nil, err
	}
	log.Printf("GetUserByID: User ID %d (%s) retrieved successfully from database", userID, user.Username)
	return &user, nil
}
