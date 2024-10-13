package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/mviner000/eyymi/eyygo/admin/models"
	"github.com/mviner000/eyymi/eyygo/shared"
)

// Configuration constants
const (
	passwordResetTimeout = 3600 // Example timeout in seconds (1 hour or 60 minutes)
)

// PasswordResetTokenGenerator handles token generation and validation.
type PasswordResetTokenGenerator struct {
	secretKey []byte
}

// NewPasswordResetTokenGenerator creates a new instance of the token generator.
func NewPasswordResetTokenGenerator() *PasswordResetTokenGenerator {
	return &PasswordResetTokenGenerator{
		secretKey: []byte(shared.GetConfig().SecretKey),
	}
}

// MakeToken generates a JWT token for the given user.
func (g *PasswordResetTokenGenerator) MakeToken(user *models.AuthUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("cannot generate token for nil user")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Second * time.Duration(passwordResetTimeout)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}

// CheckToken verifies the validity of the JWT token for the given user.
func (g *PasswordResetTokenGenerator) CheckToken(user *models.AuthUser, tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.secretKey, nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(user.ID) != claims["user_id"].(float64) {
			return false
		}
		if user.Email != claims["email"].(string) {
			return false
		}
		return true
	}

	return false
}

var DefaultTokenGenerator = NewPasswordResetTokenGenerator()
