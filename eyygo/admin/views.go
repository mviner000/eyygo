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
var tokenGenerator *auth.PasswordResetTokenGenerator

func init() {
	var err error
	dbURL := config.GetDatabaseURL(&project_name.AppSettings)
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the SQLite3 database")

	// Initialize the database connection in auth package
	auth.InitDB(db)

	// Initialize the token generator
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

	// Fetch the user from the auth service
	authUser, err := auth.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	// Password validation
	match, err := config.CheckPasswordHash(password, authUser.Password)
	if err != nil || !match {
		log.Printf("Authentication failed for user: %s", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	// Map auth.User to admin.User for additional operations
	user := &User{
		ID:        authUser.ID,
		Username:  authUser.Username,
		Email:     authUser.Email,
		Password:  authUser.Password,
		LastLogin: authUser.LastLogin,
	}

	// Generate token
	token, err := tokenGenerator.MakeToken(user.ToAuthUser())
	if err != nil {
		log.Printf("Token generation failed for user %s: %v", username, err)
		return http.HttpResponseServerError(fiber.Map{
			"error": "Failed to generate authentication token",
		}, nil).Render(c)
	}
	log.Printf("Token generated successfully for user %s", username)

	// Update last_login in the database
	if err := auth.UpdateLastLogin(user.ID); err != nil {
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

	// Redirect to the dashboard after successful login
	return http.HttpResponseRedirect("/admin/dashboard", false).Render(c)
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
