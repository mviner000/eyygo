package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mviner000/eyygo/models"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	DB          *gorm.DB
	JWTSecret   []byte
	TokenExpiry time.Duration
}

func NewAuthHandler(db *gorm.DB, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		DB:          db,
		JWTSecret:   jwtSecret,
		TokenExpiry: time.Hour * 24, // 24 hours
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := h.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !user.CheckPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Create token
	claims := jwt.MapClaims{
		"id":           user.ID,
		"username":     user.Username,
		"is_superuser": user.IsSuperUser,
		"is_staff":     user.IsStaff,
		"exp":          time.Now().Add(h.TokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := token.SignedString(h.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	h.DB.Save(&user)

	return c.JSON(fiber.Map{
		"token": t,
		"user": fiber.Map{
			"id":           user.ID,
			"username":     user.Username,
			"is_superuser": user.IsSuperUser,
			"is_staff":     user.IsStaff,
		},
	})
}

func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.JSON(fiber.Map{
		"user":  claims,
		"valid": true,
	})
}
