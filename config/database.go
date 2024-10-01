package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mviner000/eyymi/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(GetDatabaseURL()), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})
	return db
}

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

func InitConfig() {
	// Set default values
	viper.SetDefault("debug", false)
	viper.SetDefault("database.engine", "sqlite3")
	viper.SetDefault("database.name", "db.sqlite3")

	// Look for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults")
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	// Unmarshal config into AppSettings
	if err := viper.Unmarshal(&AppSettings); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	// Override with environment variables if present
	if os.Getenv("DEBUG") != "" {
		AppSettings.Debug = os.Getenv("DEBUG") == "true"
	}
	if os.Getenv("DB_ENGINE") != "" {
		AppSettings.Database.Engine = os.Getenv("DB_ENGINE")
	}
	if os.Getenv("DB_NAME") != "" {
		AppSettings.Database.Name = os.Getenv("DB_NAME")
	}

	// Ensure the database file exists (for SQLite)
	if err := EnsureDatabaseExists(); err != nil {
		log.Fatalf("Failed to ensure database exists: %v", err)
	}

	if AppSettings.Debug {
		log.Println("Configuration initialized")
	}
}
