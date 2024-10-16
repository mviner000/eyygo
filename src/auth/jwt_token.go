package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	models "github.com/mviner000/eyygo/src/admin/models"
	"github.com/mviner000/eyygo/src/config"
	"github.com/mviner000/eyygo/src/shared"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateTokenPair(userID uint, username string) (string, string, error) {
	secretKey := []byte(shared.GetSecretKey())

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // 15 minutes expiration
		"type":     "access",
	})
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
		"type":     "refresh",
	})
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func GenerateTokenPairHandler(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Retrieve user from database
	db := config.GetDB()
	var user models.AuthUser
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	// Verify password using the new function
	match, err := VerifyPassword(user.Password, body.Password) // Update to your new function
	if err != nil || !match {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User account is inactive",
		})
	}

	// Generate token pair
	accessToken, refreshToken, err := GenerateTokenPair(user.ID, user.Username)
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

func VerifyPassword(hashedPassword, password string) (bool, error) {
	secretKey := []byte(shared.GetSecretKey())
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), append([]byte(password), secretKey...))
	return err == nil, err
}

// Function to hash password (for use during user creation or password change)
func HashPassword(password string) (string, error) {
	secretKey := []byte(shared.GetSecretKey())
	hashedPassword, err := bcrypt.GenerateFromPassword(append([]byte(password), secretKey...), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
