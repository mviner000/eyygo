// project_name/middleware/jwt.go
package notes

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

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

		// Log the Authorization header for debugging
		fmt.Printf("Authorization Header: %s\n", authHeader)

		// Check if the header is empty or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		// Extract the token and log it
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("Extracted Token: %s\n", tokenString)

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Log the secret key being used
			log.Printf("Secret Key Used to Verify: %s\n", JwtSecret)

			return JwtSecret, nil
		})

		if err != nil {
			fmt.Printf("Token Parsing Error: %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if username, ok := claims["username"].(string); ok && username == "test_username" {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token or username",
		})
	}
}
