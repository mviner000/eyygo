package decorators

import (
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/mviner000/eyymi/config"
)

var sensitiveParamsRegex = regexp.MustCompile(`(?i)(pass|secret|token|key|api|pw|password)`)

// SensitivePostParameters masks sensitive POST parameters in error reports
func SensitivePostParameters() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.AppSettings.Debug {
			originalBody := c.Body()
			maskedBody := sensitiveParamsRegex.ReplaceAllString(string(originalBody), "***")
			c.Request().SetBody([]byte(maskedBody))

			err := c.Next()

			c.Request().SetBody(originalBody)
			return err
		}
		return c.Next()
	}
}

// DebugLogger logs debug information for each request
func DebugLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.AppSettings.Debug {
			log.Printf("DEBUG: %s %s", c.Method(), c.Path())
		}
		return c.Next()
	}
}
