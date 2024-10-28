// handlers/admin.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mviner000/eyygo/models"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

// ListUsers returns all users (with pagination)
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	var users []models.User
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	query := h.DB.Model(&models.User{})

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// CreateUser creates a new user
func (h *AdminHandler) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not parse user data",
		})
	}

	if err := h.DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	return c.JSON(user)
}

// UpdateUser updates user details
func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	user := new(models.User)

	if err := h.DB.First(user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not parse user data",
		})
	}

	h.DB.Save(user)
	return c.JSON(user)
}

// DeleteUser soft deletes a user
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	result := h.DB.Delete(&models.User{}, userID)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AdminMiddleware checks if user is admin
func AdminMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if !claims["is_superuser"].(bool) && !claims["is_staff"].(bool) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: Admin privileges required",
		})
	}

	return c.Next()
}

// SuperUserMiddleware checks if user is superuser
func SuperUserMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if !claims["is_superuser"].(bool) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: Superuser privileges required",
		})
	}

	return c.Next()
}
