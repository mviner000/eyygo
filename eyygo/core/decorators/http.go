package decorators

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireHTTPS redirects HTTP requests to HTTPS
func RequireHTTPS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Protocol() != "https" {
			return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
		}
		return c.Next()
	}
}

// MethodNotAllowed returns a 405 Method Not Allowed response
func MethodNotAllowed(allowedMethods ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Allow", strings.Join(allowedMethods, ", "))
		return c.SendStatus(http.StatusMethodNotAllowed)
	}
}
