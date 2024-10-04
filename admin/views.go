// admin/views.go
package admin

import (
	"github.com/gofiber/fiber/v2"
)

type UserView struct{}

// Index handles the listing of users
func (u *UserView) Index(c *fiber.Ctx) error {
	users := GetAllUsers()
	return c.Render("admin/templates/users_list", fiber.Map{
		"users": users,
	})
}

// Create renders the form for creating a new user
func (u *UserView) Create(c *fiber.Ctx) error {
	return c.Render("admin/templates/user_create", fiber.Map{})
}

// Store handles the form submission for creating a new user
func (u *UserView) Store(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")

	// Here you would typically save the user to the database
	// For now, we'll just print the values
	println("New user:", username, email)

	return c.Redirect("/admin/users")
}
