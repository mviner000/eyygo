package admin

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/auth"
	"github.com/mviner000/eyymi/eyygo/http"
	"github.com/mviner000/eyymi/project_name"
)

var store = session.New()
var db *sql.DB
var tokenGenerator *auth.PasswordResetTokenGenerator

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

	tokenGenerator = auth.NewPasswordResetTokenGenerator()
}

type User struct {
	ID            int
	Username      string
	Email         string
	Password      string
	LastLogin     time.Time
	IsSuperuser   bool
	FirstName     string
	LastName      string
	IsStaff       bool
	IsActive      bool
	DateJoined    time.Time
	GroupsID      sql.NullInt64
	PermissionsID sql.NullInt64
}

// Add this function to convert admin.User to auth.User
func (u *User) ToAuthUser() *auth.User {
	return &auth.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		LastLogin: u.LastLogin,
	}
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

	match, err := config.CheckPasswordHash(password, user.Password)
	if err != nil || !match {
		log.Printf("Authentication failed for user: %s", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	// Generate token
	authUser := user.ToAuthUser()
	token, err := tokenGenerator.MakeToken(authUser)
	if err != nil {
		log.Printf("Token generation failed for user %s: %v", username, err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to generate authentication token",
		}, nil).Render(c)
	}
	log.Printf("Token generated successfully for user %s", username)

	log.Printf("DEBUG: Generated token for user %s: %s", username, token)

	// Update last_login in the database
	if err := updateLastLogin(user.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", username, err)
		// Continue with login process even if update fails
	} else {
		log.Printf("Last login updated successfully for user %s", username)
	}

	// Create session
	sess, err := store.Get(c)
	if err != nil {
		log.Printf("Session creation failed for user %s: %v", username, err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to create session",
		}, nil).Render(c)
	}

	sess.Set("user_id", user.ID)
	sess.Set("auth_token", token)
	if err := sess.Save(); err != nil {
		log.Printf("Session save failed for user %s: %v", username, err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to save session",
		}, nil).Render(c)
	}

	log.Printf("User authenticated successfully: %s", user.Username)
	return http.HttpResponseRedirect("/admin/dashboard", false).Render(c)
}

func getUserByUsername(username string) (*User, error) {
	var user User
	var lastLogin sql.NullTime
	query := "SELECT id, username, email, password, last_login FROM auth_user WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &lastLogin)
	if err != nil {
		log.Printf("Error retrieving user %s from database: %v", username, err)
		return nil, err
	}
	if lastLogin.Valid {
		user.LastLogin = lastLogin.Time
	} else {
		user.LastLogin = time.Time{} // Set to zero time if NULL
	}
	log.Printf("User %s retrieved successfully from database", username)
	return &user, nil
}

func updateLastLogin(userID int) error {
	_, err := db.Exec("UPDATE auth_user SET last_login = ? WHERE id = ?", time.Now(), userID)
	if err != nil {
		log.Printf("Error updating last_login for user ID %d: %v", userID, err)
	} else {
		log.Printf("Last login updated successfully for user ID %d", userID)
	}
	return err
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
