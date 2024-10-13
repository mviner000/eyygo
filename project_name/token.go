// project_name/token/token.go
package project_name

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	middleware "github.com/mviner000/eyymi/project_name/notes"
)

func GenerateTokenPair(username string) (string, string, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // 15 minutes expiration
	})
	accessTokenString, err := accessToken.SignedString(middleware.JwtSecret)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
	})
	refreshTokenString, err := refreshToken.SignedString(middleware.JwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
