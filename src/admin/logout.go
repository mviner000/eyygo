package admin

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyygo/src/auth"
	"github.com/mviner000/eyygo/src/http"
)

func Logout(c *fiber.Ctx) error {
	log.Println("Logout function called.")

	// Use the constant for the session cookie name
	sessionID := c.Cookies(auth.SessionCookieName)
	if sessionID != "" {
		log.Printf("Session ID: %s", sessionID)
		err := auth.DeleteSessionFromDB(sessionID)
		if err != nil {
			log.Printf("Error deleting session from DB: %v", err)
		}
	} else {
		log.Println("No session ID found.")
	}

	// Clear the session cookie using the utility function
	auth.DeleteSessionCookie(c)

	auth.SetSessionCookie(c, sessionID, time.Now(), 0)

	log.Println("Session cookie cleared.")

	// Redirect to the login page
	return c.SendString(http.WindowReload("/admin/login"))
}
