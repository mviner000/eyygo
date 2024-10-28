// views/views.go
package views

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyygo/models"
	"gorm.io/gorm"
)

type ViewHandler struct {
	DB *gorm.DB
}

func NewViewHandler(db *gorm.DB) *ViewHandler {
	return &ViewHandler{DB: db}
}

// LoginPage renders the login page
func (h *ViewHandler) LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title":       "Login",
		"CurrentYear": time.Now().Year(),
	}, "layouts/auth")
}

// ValidateUsername handles username validation
func (h *ViewHandler) ValidateUsername(c *fiber.Ctx) error {
	username := c.FormValue("username")
	if username == "" {
		return c.Status(422).SendString(`
            <div class="text-red-500 text-sm mt-1">
                Username is required
            </div>
        `)
	}
	return c.SendString("")
}

// ValidatePassword handles password validation
func (h *ViewHandler) ValidatePassword(c *fiber.Ctx) error {
	password := c.FormValue("password")
	if password == "" {
		return c.Status(422).SendString(`
            <div class="text-red-500 text-sm mt-1">
                Password is required
            </div>
        `)
	}
	return c.SendString("")
}

// ValidateLogin handles the login validation
func (h *ViewHandler) ValidateLogin(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		return c.Status(422).SendString(`
            <div class="text-red-500 text-sm mt-1">
                Both username and password are required
            </div>
        `)
	}

	// Check user credentials
	var user models.User
	if err := h.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(401).SendString(`
            <div class="text-red-500 text-sm mt-1">
                Invalid credentials
            </div>
        `)
	}

	if !user.CheckPassword(password) {
		return c.Status(401).SendString(`
            <div class="text-red-500 text-sm mt-1">
                Invalid credentials
            </div>
        `)
	}

	// Set session/token here
	// Add success message header
	c.Response().Header.Add("HX-Trigger", `{"showMessage": "Login successful"}`)
	c.Response().Header.Add("HX-Redirect", "/dashboard")
	return c.SendString("")
}

// Dashboard renders the main dashboard
func (h *ViewHandler) Dashboard(c *fiber.Ctx) error {
	var userCount int64
	var noteCount int64

	h.DB.Model(&models.User{}).Count(&userCount)
	h.DB.Model(&models.Note{}).Count(&noteCount)

	return c.Render("dashboard", fiber.Map{
		"Title":     "Dashboard",
		"UserCount": userCount,
		"NoteCount": noteCount,
	}, "layouts/main")
}

// UsersList handles the HTMX request for users list
func (h *ViewHandler) UsersList(c *fiber.Ctx) error {
	var users []models.User
	result := h.DB.Find(&users)
	if result.Error != nil {
		return c.Status(500).SendString("Error loading users")
	}

	return c.Render("users-list", fiber.Map{
		"Users": users,
	})
}

// NotesList handles the HTMX request for notes list
func (h *ViewHandler) NotesList(c *fiber.Ctx) error {
	var notes []models.Note
	result := h.DB.Find(&notes)
	if result.Error != nil {
		return c.Status(500).SendString("Error loading notes")
	}

	return c.Render("notes-list", fiber.Map{
		"Notes": notes,
	})
}

// LogoutPage handles user logout
func (h *ViewHandler) LogoutPage(c *fiber.Ctx) error {
	// Clear session/token here
	c.Response().Header.Add("HX-Redirect", "/login")
	return c.SendString("")
}
