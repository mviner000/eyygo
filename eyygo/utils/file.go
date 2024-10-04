package utils

import (
	"fmt"
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
