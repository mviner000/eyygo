package config

import (
	"encoding/base64"
	"fmt"

	"github.com/mviner000/eyymi/eyygo/shared" // Import your project's package
	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12
)

func HashPassword(password string) (string, error) {
	secretKey := shared.GetConfig().SecretKey
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+secretKey), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	secretKey := shared.GetConfig().SecretKey
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("error decoding hash: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(decodedHash, []byte(password+secretKey))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error comparing password hash: %w", err)
	}

	return true, nil
}
