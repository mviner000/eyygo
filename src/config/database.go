// File: config/database.go

package config

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/mviner000/eyymi/src/shared"
	"github.com/mviner000/eyymi/src/utils"
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
	sharedConfig := shared.GetConfig()
	db := sharedConfig.Database
	var dbURL string
	switch db.Engine {
	case "sqlite3":
		dbURL = db.Name
		if ext := filepath.Ext(dbURL); ext != ".db" && ext != ".sqlite3" {
			dbURL += ".db"
		}
		if !filepath.IsAbs(dbURL) {
			dbURL = filepath.Join(utils.GetProjectRoot(sharedConfig.Debug), dbURL)
		}
	case "mysql":
		dbURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.User, db.Password, db.Host, db.Port, db.Name)
	case "postgres":
		dbURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			db.Host, db.Port, db.User, db.Name, db.Password)
	default:
		log.Printf("Unsupported database engine: %s, falling back to SQLite", db.Engine)
		dbPath, err := filepath.Abs("db.sqlite3")
		if err != nil {
			log.Printf("Error getting absolute path for default database: %v", err)
			dbPath = "db.sqlite3"
		}
		dbURL = dbPath
	}
	if sharedConfig.Debug {
		log.Printf("Database URL: %s", dbURL)
	}
	return dbURL
}

func EnsureDatabaseExists() error {
	sharedConfig := shared.GetConfig()
	if sharedConfig.Database.Engine == "sqlite3" {
		dbPath := filepath.Join(utils.GetProjectRoot(sharedConfig.Debug), sharedConfig.Database.Name)
		return utils.EnsureFileExists(dbPath)
	}
	return nil
}

// InitDatabaseConfig initializes the database configuration in the shared config
func InitDatabaseConfig(dbConfig shared.DatabaseConfig, debug bool) {
	shared.SetDatabaseConfig(dbConfig)
	shared.SetDebug(debug)
}
