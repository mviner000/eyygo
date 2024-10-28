package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server settings
	ServerPort string
	ServerHost string

	// Database settings
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		// Server settings
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		// Database settings
		DBDriver:   getEnv("DB_DRIVER", "sqlite"), // sqlite, mysql, postgresql
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"), // default MySQL port
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "test_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
