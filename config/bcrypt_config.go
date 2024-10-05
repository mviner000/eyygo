package config

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 12
	BcryptKey  = "your-secret-key-here" // Replace with a secure, random string
)

func HashPassword(password string) (string, error) {
	fmt.Println("Starting password hashing")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+BcryptKey), BcryptCost)
	if err != nil {
		fmt.Println("Error occurred during password hashing:", err)
	} else {
		fmt.Println("Password hashing completed successfully")
	}
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	fmt.Println("Starting password hash verification")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+BcryptKey))
	if err != nil {
		fmt.Println("Password hash verification failed:", err)
	} else {
		fmt.Println("Password hash verification succeeded")
	}
	return err == nil
}
