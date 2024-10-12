package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/shared"
)

// Configuration constants
const (
	keySalt              = "dfdsfsdfdrato"
	passwordResetTimeout = 3600 // Example timeout in seconds
)

// PasswordResetTokenGenerator handles token generation and validation.
type PasswordResetTokenGenerator struct {
	secret string
}

// NewPasswordResetTokenGenerator creates a new instance of the token generator.
func NewPasswordResetTokenGenerator() *PasswordResetTokenGenerator {
	return &PasswordResetTokenGenerator{
		secret: shared.GetConfig().SecretKey,
	}
}

// MakeToken generates a token for the given user.
func (g *PasswordResetTokenGenerator) MakeToken(user *models.AuthUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("cannot generate token for nil user")
	}
	timestamp := time.Now().Unix()
	token := g.makeTokenWithTimestamp(user, timestamp)
	return token, nil
}

// CheckToken verifies the validity of the token for the given user.
func (g *PasswordResetTokenGenerator) CheckToken(user *models.AuthUser, token string) bool {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return false
	}
	tsStr := parts[0]
	ts, err := base36ToInt(tsStr)
	if err != nil {
		return false
	}
	if time.Now().Unix()-int64(ts) > passwordResetTimeout {
		return false
	}
	expectedToken := g.makeTokenWithTimestamp(user, int64(ts))

	log.Printf("Token check: User ID: %d, Timestamp: %d, Expected Token: %s, Received Token: %s", user.ID, ts, expectedToken, token)
	return constantTimeCompare(expectedToken, token)
}

// makeTokenWithTimestamp creates a token using a timestamp.
func (g *PasswordResetTokenGenerator) makeTokenWithTimestamp(user *models.AuthUser, timestamp int64) string {
	tsB36 := intToBase36(int(timestamp))
	hashValue := g.makeHashValue(user, timestamp)
	hashString := g.saltedHMAC(keySalt, hashValue, g.secret)
	token := fmt.Sprintf("%s-%s", tsB36, hashString[:len(hashString)/2])
	return token
}

// makeHashValue generates a hash value based on user details and timestamp.
func (g *PasswordResetTokenGenerator) makeHashValue(user *models.AuthUser, timestamp int64) string {
	return fmt.Sprintf("%d%s%d%s", user.ID, user.Password, timestamp, user.Email)
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
