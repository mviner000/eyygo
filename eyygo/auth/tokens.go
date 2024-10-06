package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mviner000/eyymi/project_name"
)

// Configuration constants
const (
	keySalt              = "dfdsfsdfdrato"
	passwordResetTimeout = 3600 // Example timeout in seconds
)

// PasswordResetTokenGenerator handles token generation and validation.
type PasswordResetTokenGenerator struct {
	secret          string
	secretFallbacks []string
}

// NewPasswordResetTokenGenerator creates a new instance of the token generator.
func NewPasswordResetTokenGenerator() *PasswordResetTokenGenerator {
	return &PasswordResetTokenGenerator{
		secret: project_name.AppSettings.SecretKey,
	}
}

func (g *PasswordResetTokenGenerator) numSeconds(t time.Time) int {
	return int(t.Unix())
}

// MakeToken generates a token for the given user.
func (g *PasswordResetTokenGenerator) MakeToken(user *User) (string, error) {
	if user == nil {
		log.Println("Error: Attempted to generate token for nil user")
		return "", fmt.Errorf("cannot generate token for nil user")
	}
	timestamp := g.numSeconds(time.Now())
	token := g.makeTokenWithTimestamp(user, timestamp, g.secret)
	log.Printf("Token generated successfully for user ID %d", user.ID)
	return token, nil
}

// CheckToken verifies the validity of the token for the given user.
func (g *PasswordResetTokenGenerator) CheckToken(user *User, token string) bool {
	if user == nil || token == "" {
		log.Println("Error: Attempted to check token with nil user or empty token")
		return false
	}

	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		log.Printf("Error: Invalid token format for user ID %d", user.ID)
		return false
	}

	tsB36 := parts[0]
	ts, err := base36ToInt(tsB36)
	if err != nil {
		log.Printf("Error: Failed to parse timestamp from token for user ID %d: %v", user.ID, err)
		return false
	}

	for _, secret := range append([]string{g.secret}, g.secretFallbacks...) {
		expectedToken := g.makeTokenWithTimestamp(user, ts, secret)
		if constantTimeCompare(expectedToken, token) {
			if time.Now().Unix()-int64(ts) <= passwordResetTimeout {
				log.Printf("Token validated successfully for user ID %d", user.ID)
				return true
			} else {
				log.Printf("Token expired for user ID %d", user.ID)
			}
		}
	}

	log.Printf("Token validation failed for user ID %d", user.ID)
	return false
}

// makeTokenWithTimestamp creates a token using a timestamp and secret.
func (g *PasswordResetTokenGenerator) makeTokenWithTimestamp(user *User, timestamp int, secret string) string {
	tsB36 := intToBase36(timestamp)
	hashString := g.saltedHMAC(keySalt, g.makeHashValue(user, timestamp), secret)
	return fmt.Sprintf("%s-%s", tsB36, hashString[:len(hashString)/2])
}

// makeHashValue generates a hash value based on user details and timestamp.
func (g *PasswordResetTokenGenerator) makeHashValue(user *User, timestamp int) string {
	loginTimestamp := ""
	if !user.LastLogin.IsZero() {
		loginTimestamp = user.LastLogin.UTC().Format(time.RFC3339)
	}
	return fmt.Sprintf("%d%s%s%d%s", user.ID, user.Password, loginTimestamp, timestamp, user.Email)
}

// saltedHMAC performs a salted HMAC operation using SHA-256.
func (g *PasswordResetTokenGenerator) saltedHMAC(keySalt, value, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(keySalt + value))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// constantTimeCompare compares two strings in a way that is resistant to timing attacks.
func constantTimeCompare(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

// base36ToInt converts a base36 encoded string to an integer.
func base36ToInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%x", &n)
	if err != nil {
		return 0, fmt.Errorf("invalid base36 string: %w", err)
	}
	return n, nil
}

// intToBase36 converts an integer to a base36 encoded string.
func intToBase36(n int) string {
	return fmt.Sprintf("%x", n)
}

// now returns the current time.
func now() time.Time {
	return time.Now()
}

var DefaultTokenGenerator = NewPasswordResetTokenGenerator()
