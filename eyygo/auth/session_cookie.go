package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const SessionCookieName = "hey"

// SetSessionCookie sets a session cookie with the provided session ID, expiration time, and max age.
func SetSessionCookie(c *fiber.Ctx, sessionID string, expireTime time.Time, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieName, // Corrected spelling
		Value:    sessionID,
		Expires:  expireTime,
		MaxAge:   maxAge, // Use the provided MaxAge
		HTTPOnly: true,
		Secure:   true, // Ensure to set this based on your environment (use false in local)
	})
}

func DeleteSessionCookie(c *fiber.Ctx) {
	c.ClearCookie(SessionCookieName)
}
