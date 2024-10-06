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

	authUser, err := auth.GetUserByUsername(username)
	if err != nil {
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	match, err := config.CheckPasswordHash(password, authUser.Password)
	if err != nil || !match {
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	tokenGenerator := auth.NewPasswordResetTokenGenerator()
	token, err := tokenGenerator.MakeToken(authUser)
	if err != nil {
		return http.HttpResponseServerError("Failed to generate authentication token", nil).Render(c)
	}

	// Update last_login in the database
	if err := auth.UpdateLastLogin(authUser.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", username, err)
	}

	// Create session
	sessionID := generateSessionID()
	expireTime := time.Now().Add(24 * time.Hour)

	// Store session in the database
	query := `INSERT INTO eyygo_session (session_key, user_id, auth_token, expire_date) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, sessionID, authUser.ID, token, expireTime)
	if err != nil {
		return http.HttpResponseServerError("Failed to create session", nil).Render(c)
	}

	// Set the session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expireTime,
		HTTPOnly: true,
		Secure:   true,
	})

	return http.HttpResponseRedirect("/admin/dashboard", false).Render(c)
}

// generateSessionID generates a new session ID
func generateSessionID() string {
	return uuid.New().String()
}

func Dashboard(c *fiber.Ctx) error {
	userID, _, err := auth.GetSessionFromDB(c)
	if err != nil {
		return http.HttpResponseRedirect("/login", false).Render(c)
	}

	user, err := auth.GetUserByID(userID)
	if err != nil {
		return http.HttpResponseServerError("Error retrieving user information", nil).Render(c)
	}

	log.Printf("User data: %+v", user)

	return http.HttpResponseHTMX(fiber.Map{
		"User": user, // No space after "User"
	}, "eyygo/admin/templates/dashboard.html", "eyygo/admin/templates/layout.html").Render(c)
}

func UserList(c *fiber.Ctx) error {
	users, err := auth.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving user list")
	}

	return c.Render("eyygo/admin/templates/user_list", fiber.Map{
		"Users": users,
	})
}

func UserCreate(c *fiber.Ctx) error {
	response := http.HttpResponseOK(fiber.Map{}, nil, "eyygo/admin/templates/user_form")
	return response.Render(c)
}

func UserStore(c *fiber.Ctx) error {
	return c.SendString("User creation logic not implemented")
}
