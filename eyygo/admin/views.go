package admin

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/mviner000/eyymi/eyygo/auth"
	"github.com/mviner000/eyymi/eyygo/config"
	"github.com/mviner000/eyymi/eyygo/http"

	models "github.com/mviner000/eyymi/project_name/posts"
)

const SiteName = "Eyygo Administration"

var store = session.New()
var tokenGenerator *auth.PasswordResetTokenGenerator

func init() {
	// Initialize the database connection in the config package
	db := config.GetDB()
	if db == nil {
		log.Fatalf("Failed to connect to database")
	}
	log.Println("Successfully connected to the database")

	// Pass *sql.DB to auth.InitDB
	auth.InitDB(db)

	// Initialize the token generator
	tokenGenerator = auth.NewPasswordResetTokenGenerator()
}

func LoginForm(c *fiber.Ctx) error {
	log.Println("LoginForm function called")
	errorMessage := c.Query("error")

	// Check if it's an HTMX request
	if c.Get("HX-Request") == "true" {
		log.Println("HTMX request detected")
		// If it's an HTMX request, just return the form content
		return http.HttpResponseHTMX(fiber.Map{
			"Error": errorMessage,
		}, "eyygo/admin/templates/login_form.html").Render(c)
	}

	log.Println("Rendering full login page")
	// Render the full page with layout
	return http.HttpResponseHTMX(fiber.Map{
		"Error":     errorMessage,
		"MetaTitle": "Login | " + SiteName,
	}, "eyygo/admin/templates/login.html", "eyygo/admin/templates/layout.html").Render(c)
}

func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	authUser, err := auth.GetUserByUsername(username)
	if err != nil {
		log.Printf("Login attempt failed for username '%s': user not found", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	match, err := config.CheckPasswordHash(password, authUser.Password)
	if err != nil || !match {
		log.Printf("Login attempt failed for user '%s': invalid password", username)
		return http.HttpResponseUnauthorized(fiber.Map{
			"error": "Invalid username or password",
		}, nil).Render(c)
	}

	token, err := tokenGenerator.MakeToken(authUser)

	if err != nil {
		log.Printf("Failed to generate authentication token for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to generate authentication token", nil).Render(c)
	}
	log.Printf("Generated token for user '%s': %s", username, token)

	// Update last_login in the database
	if err := auth.UpdateLastLogin(authUser.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", username, err)
	}

	// Create session
	sessionID := generateSessionID()
	expireTime := time.Now().Add(24 * time.Hour)

	// Store session in the database
	if authUser.ID < 0 {
		log.Printf("Invalid user ID: %d", authUser.ID)
		return http.HttpResponseServerError("Invalid user ID", nil).Render(c)
	}

	session := models.Session{
		SessionKey: sessionID,
		UserID:     uint(authUser.ID), // Convert int to uint
		AuthToken:  token,
		ExpireDate: expireTime,
	}

	if err := config.GetDB().Create(&session).Error; err != nil {
		log.Printf("Failed to create session for user '%s': %v", username, err)
		return http.HttpResponseServerError("Failed to create session", nil).Render(c)
	}
	log.Printf("Created session for user '%s': sessionID=%s, token=%s", username, sessionID, token)

	// Set the session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "hey_sesion",
		Value:    sessionID,
		Expires:  expireTime,
		HTTPOnly: true,
		Secure:   true,
	})

	log.Printf("Login successful for user '%s'. Session ID: %s, Auth Token: %s", username, sessionID, token)

	return c.SendString(http.WindowReload("/admin/dashboard"))
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

	// Convert int to uint
	if userID < 0 {
		return http.HttpResponseServerError("Invalid user ID", nil).Render(c)
	}
	uintUserID := uint(userID)

	user, err := auth.GetUserByID(uintUserID)
	if err != nil {
		return http.HttpResponseServerError("Error retrieving user information", nil).Render(c)
	}

	log.Printf("User data: %+v", user)

	return http.HttpResponseHTMX(fiber.Map{
		"User":      user,
		"MetaTitle": "Dashboard | " + SiteName,
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

func Logout(c *fiber.Ctx) error {
	log.Println("Logout function called.")

	sessionID := c.Cookies("hey_sesion")
	if sessionID != "" {
		log.Printf("Session ID: %s", sessionID)
		err := auth.DeleteSessionFromDB(sessionID)
		if err != nil {
			log.Printf("Error deleting session from DB: %v", err)
		}
	} else {
		log.Println("No session ID found.")
	}

	// Clear the session cookie
	c.ClearCookie("hey_sesion")
	log.Println("Session cookie cleared.")

	return c.SendString(http.WindowReload("/admin/login"))
}
