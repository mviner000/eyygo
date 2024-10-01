package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mviner000/eyymi/utils"
)

func GetDatabaseURL() string {
	db := AppSettings.Database
	var dbURL string
	switch db.Engine {
	case "sqlite3":
		cwd, err := os.Getwd()
		if err != nil {
			if AppSettings.Debug {
				log.Printf("Error getting current working directory: %v", err)
			}
			cwd = "."
		}
		dbPath := filepath.Join(cwd, db.Name)
		dbURL = dbPath // Ent expects the file path for SQLite, not a URL
	// ... (keep other database cases)
	default:
		if AppSettings.Debug {
			log.Printf("Unsupported database engine: %s, falling back to SQLite", db.Engine)
		}
		cwd, err := os.Getwd()
		if err != nil {
			if AppSettings.Debug {
				log.Printf("Error getting current working directory: %v", err)
			}
			cwd = "."
		}
		dbPath := filepath.Join(cwd, "db.sqlite3")
		dbURL = dbPath
	}
	if AppSettings.Debug {
		log.Printf("Database URL: %s", dbURL)
	}
	return dbURL
}

func EnsureDatabaseExists() error {
	if AppSettings.Database.Engine == "sqlite3" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %v", err)
		}
		dbPath := filepath.Join(cwd, AppSettings.Database.Name)
		return utils.EnsureFileExists(dbPath)
	}
	return nil
}
