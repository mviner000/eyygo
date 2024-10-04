// admin/views.go
package admin

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/config"
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

	// Open the database connection
	db, err := sql.Open("sqlite3", config.GetDatabaseURL())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to connect to database")
	}
	defer db.Close()

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", username, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert user")
	}

	return c.Redirect("/admin/users")
}
