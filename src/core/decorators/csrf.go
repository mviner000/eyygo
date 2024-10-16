package decorators

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	conf "github.com/mviner000/eyymi/src"
	"github.com/mviner000/eyymi/src/config"
)

var (
	csrfTokens     = make(map[string]time.Time)
	csrfTokenMutex sync.RWMutex
)

// CSRFConfig holds the configuration for CSRF protection
type CSRFConfig struct {
	TokenLength   int
	CookieName    string
	HeaderName    string
	Expiration    time.Duration
	SecureCookie  bool
	ExemptMethods []string
}

// DefaultCSRFConfig returns the default CSRF configuration
func DefaultCSRFConfig(cfg config.Config) CSRFConfig {
	return CSRFConfig{
		TokenLength:   32,
		CookieName:    "csrf_token",
		HeaderName:    "X-CSRF-Token",
		Expiration:    1 * time.Hour,
		SecureCookie:  !cfg.IsDebug(), // Use the IsDebug method to determine if it's a development environment
		ExemptMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// generateToken creates a new CSRF token
func generateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CSRFMiddleware provides CSRF protection
func CSRFMiddleware(cfg config.Config, customCfg ...CSRFConfig) fiber.Handler {
	// Use default config
	csrfConfig := DefaultCSRFConfig(cfg)
	if len(customCfg) > 0 {
		csrfConfig = customCfg[0]
	}

	return func(c *fiber.Ctx) error {
		// Check if the method is exempt
		for _, method := range csrfConfig.ExemptMethods {
			if c.Method() == method {
				return c.Next()
			}
		}

		// Get the token from the request
		token := c.Cookies(csrfConfig.CookieName)
		if token == "" {
			token = c.Get(csrfConfig.HeaderName)
		}

		csrfTokenMutex.RLock()
		tokenTime, exists := csrfTokens[token]
		csrfTokenMutex.RUnlock()

		// If token doesn't exist or has expired, generate a new one
		if !exists || time.Now().After(tokenTime) {
			var err error
			token, err = generateToken(csrfConfig.TokenLength)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to generate CSRF token",
				})
			}

			csrfTokenMutex.Lock()
			csrfTokens[token] = time.Now().Add(csrfConfig.Expiration)
			csrfTokenMutex.Unlock()

			c.Cookie(&fiber.Cookie{
				Name:     csrfConfig.CookieName,
				Value:    token,
				Expires:  time.Now().Add(csrfConfig.Expiration),
				Secure:   csrfConfig.SecureCookie,
				HTTPOnly: true,
				SameSite: "Strict",
			})
		}

		// For non-GET requests, validate the token
		if c.Method() != "GET" {
			requestToken := c.Get(csrfConfig.HeaderName)
			if requestToken == "" || requestToken != token {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Invalid CSRF token",
				})
			}
		}

		// Set the token in the response headers
		c.Set(csrfConfig.HeaderName, token)

		return c.Next()
	}
}

// CSRFToken returns the current CSRF token
func CSRFToken(c *fiber.Ctx) string {
	cfg := conf.GetSettings()
	return c.Get(DefaultCSRFConfig(cfg).HeaderName)
}

// CleanupCSRFTokens removes expired tokens
func CleanupCSRFTokens() {
	csrfTokenMutex.Lock()
	defer csrfTokenMutex.Unlock()

	for token, expiry := range csrfTokens {
		if time.Now().After(expiry) {
			delete(csrfTokens, token)
		}
	}
}
