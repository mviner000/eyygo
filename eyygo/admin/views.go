package admin

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mviner000/eyymi/config"
	"github.com/mviner000/eyymi/eyygo/auth"
	"github.com/mviner000/eyymi/eyygo/http"
	"github.com/mviner000/eyymi/project_name"
)

var db *sql.DB

func InitDB(database *sql.DB) {
	db = database
}

var store = session.New()
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

// Convert admin.User to auth.User
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
		return http.HttpResponseServerError("Failed to generate authentication token", map[string]string{
			"error": "Failed to generate authentication token",
		}).Render(c)
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
	sessionID := generateSessionID()             // Implement this function to generate a unique session ID
	expireTime := time.Now().Add(24 * time.Hour) // Set session expiry to 24 hours from now

	// Store session in the database
	query := `INSERT INTO eyygo_session (session_key, user_id, auth_token, expire_date) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, sessionID, user.ID, token, expireTime)
	if err != nil {
		log.Printf("Failed to store session in database for user %s: %v", username, err)
		return http.HttpResponseServerError("Failed to create session", map[string]string{
			"error": "Failed to create session",
		}).Render(c)
	}
	log.Printf("Session created successfully for user %s", username)

	// Set the session cookie in the response
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expireTime,
		HTTPOnly: true,
		Secure:   true,
	})

	response := http.HttpResponseOK(fiber.Map{
		"message": "Login successful",
	}, nil, "eyygo/admin/templates/dashboard")
	return response.Render(c)
}

// generateSessionID generates a new session ID
func generateSessionID() string {
	return uuid.New().String()
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
