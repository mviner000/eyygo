package admin

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/eyygo/auth"
	"github.com/mviner000/eyymi/eyygo/http"
)

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
