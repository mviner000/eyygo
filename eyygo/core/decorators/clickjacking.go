package decorators

import (
	"github.com/gofiber/fiber/v2"
)

// XFrameOptions represents the possible values for the X-Frame-Options header
type XFrameOptions string

const (
	DENY       XFrameOptions = "DENY"
	SAMEORIGIN XFrameOptions = "SAMEORIGIN"
)

// XFrameOptionsMiddleware adds X-Frame-Options header to the response
func XFrameOptionsMiddleware(option XFrameOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", string(option))
		return c.Next()
	}
}

// XFrameOptionsDeny is a convenience middleware that sets X-Frame-Options to DENY
func XFrameOptionsDeny(c *fiber.Ctx) error {
	return XFrameOptionsMiddleware(DENY)(c)
}

// XFrameOptionsSameOrigin is a convenience middleware that sets X-Frame-Options to SAMEORIGIN
func XFrameOptionsSameOrigin(c *fiber.Ctx) error {
	return XFrameOptionsMiddleware(SAMEORIGIN)(c)
}

// ContentSecurityPolicyMiddleware adds Content-Security-Policy header to the response
func ContentSecurityPolicyMiddleware(policy string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", policy)
		return c.Next()
	}
}

// FrameOptionsPolicyMiddleware adds Frame-Options header to the response (for older browsers)
func FrameOptionsPolicyMiddleware(option XFrameOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Frame-Options", string(option))
		return c.Next()
	}
}
