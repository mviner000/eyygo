package admin

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/http"
	"github.com/mviner000/eyymi/project_name"
)

var store = session.New()
var db *sql.DB

func init() {
	var err error
	dbURL := config.GetDatabaseURL(&project_name.AppSettings)

	db, err = sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the SQLite3 database")
}

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

func LoginForm(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/login")
	return response.Render(c)
}

func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	log.Printf("Login attempt: username=%s", username)

	if username == "" || password == "" {
		log.Println("Validation failed: Username and password are required")
		return http.HttpResponseBadRequest(fiber.Map{
			"error": "Username and password are required",
		}, nil).Render(c)
	}

	user, err := getUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	log.Printf("Retrieved user from database: %+v", user)

	if !config.CheckPasswordHash(password, user.Password) {
		log.Printf("Password comparison failed for user: %s", user.Username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	log.Printf("User authenticated: %s", user.Username)

	session, err := store.Get(c)
	if err != nil {
		log.Printf("Session creation failed: %v", err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to create session",
		}, nil).Render(c)
	}

	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Printf("Session save failed: %v", err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to save session",
		}, nil).Render(c)
	}

	log.Printf("Session created for user_id: %d", user.ID)

	return http.HttpResponseRedirect("/admin/dashboard", false).Render(c)
}

func getUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, email, password FROM auth_user WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &user, nil
}

func Dashboard(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/dashboard")
	return response.Render(c)
}

func UserIndex(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/users_list")
	return response.Render(c)
}

func UserCreate(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/user_form")
	return response.Render(c)
}

func UserStore(c *fiber.Ctx) error {
	return c.SendString("User creation logic not implemented")
}
