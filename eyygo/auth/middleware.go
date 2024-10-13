package auth

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/shared"
)

var JwtSecret []byte

func InitJWTSecret() {
	JwtSecret = []byte(shared.GetSecretKey())
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the header is empty or doesn't start with "Bearer "
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header",
			})
		}

		// Extract the token
		tokenString := authHeader[7:]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Check if the token is valid and extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check token expiration
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "Token expired",
					})
				}
			}

			// Extract user ID from claims
			userID, ok := claims["user_id"].(float64)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims",
				})
			}

			// Fetch user from database
			db := config.GetDB()
			var user models.AuthUser
			if err := db.First(&user, uint(userID)).Error; err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User not found",
				})
			}

			// Check if user is active
			if !user.IsActive {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User is inactive",
				})
			}

			// Store user in context for later use
			c.Locals("user", user)

			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
}

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
