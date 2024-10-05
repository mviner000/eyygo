package admin

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/app_name"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/http"
)

type UserView struct{}

// Index handles the listing of users
func (u *UserView) Index(c *fiber.Ctx) error {
	users := GetAllUsers()
	response := http.HttpResponseOK(fiber.Map{
		"users": users,
	}, nil, "eyygo/admin/templates/users_list") // Provide the template name
	return response.Render(c)
}

// Create renders the form for creating a new user
func (u *UserView) Create(c *fiber.Ctx) error {
	return http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/user_create").Render(c) // Provide the template name
}

// Store handles the form submission for creating a new user
func (u *UserView) Store(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")

	// Use the app_name.AppSettings to get the database URL
	dbURL := config.GetDatabaseURL(&app_name.AppSettings)

	// Open the database connection
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return http.HttpResponseServerError(err.Error(), nil).Render(c)
	}
	defer db.Close()

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", username, email)
	if err != nil {
		return http.HttpResponseServerError(err.Error(), nil).Render(c)
	}

	// Redirect to the users list page
	return http.HttpResponseRedirect("/admin/users", false).Render(c)
}
