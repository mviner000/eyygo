package project_name

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mviner000/eyymi/eyygo/shared"
)

func GenerateTokenPair(username string) (string, string, error) {
	secretKey := []byte(shared.GetSecretKey())

	// Log the secret key being used to sign the tokens
	log.Printf("Secret Key Used to Sign: %s\n", secretKey)

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // 15 minutes expiration
	})
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Log the generated access token for debugging purposes
	log.Printf("Generated Access Token: %s\n", accessTokenString)

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiration
	})
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Log the generated refresh token for debugging purposes
	log.Printf("Generated Refresh Token: %s\n", refreshTokenString)

	return accessTokenString, refreshTokenString, nil
}
