package utils

import (
	"fmt"
	"log"
	"os"
)

func EnsureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		file.Close()
	}
	return nil
}

// GetProjectRoot returns the current working directory.
// If there's an error getting the current working directory, it returns "." and logs the error if debug is true.
func GetProjectRoot(debug bool) string {
	cwd, err := os.Getwd()
	if err != nil {
		if debug {
			log.Printf("Error getting current working directory: %v", err)
		}
		cwd = "."
	}
	return cwd
}
