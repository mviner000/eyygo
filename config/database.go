package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mviner000/eyymi/eyygo/utils"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB(cfg Config) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(GetDatabaseURL(cfg)), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})
	return db
}

func GetProjectRoot(cfg Config) string {
	cwd, err := os.Getwd()
	if err != nil {
		if cfg.IsDebug() {
			log.Printf("Error getting current working directory: %v", err)
		}
		cwd = "."
	}
	return cwd
}

func GetDatabaseURL(cfg Config) string {
	db := cfg.GetDatabaseConfig()
	var dbURL string
	switch db.Engine {
	case "sqlite3":
		// For SQLite3, use the DB_NAME as the file path
		dbURL = db.Name
		// If DB_NAME doesn't have a .db or .sqlite3 extension, add .db
		if ext := filepath.Ext(dbURL); ext != ".db" && ext != ".sqlite3" {
			dbURL += ".db"
		}
		// If it's not an absolute path, make it relative to the current directory
		if !filepath.IsAbs(dbURL) {
			dbURL = filepath.Join(".", dbURL)
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

	if cfg.IsDebug() {
		log.Printf("Database URL: %s", dbURL)
	}
	return dbURL
}

func GetDatabaseURLForDbmate(cfg Config) string {
	db := cfg.GetDatabaseConfig()
	var dbURL string
	switch db.Engine {
	case "sqlite3":
		cwd, err := os.Getwd()
		if err != nil {
			if cfg.IsDebug() {
				log.Printf("Error getting current working directory: %v", err)
			}
			cwd = "."
		}
		dbPath := filepath.Join(cwd, db.Name)
		dbURL = fmt.Sprintf("sqlite3://%s", dbPath)
	default:
		if cfg.IsDebug() {
			log.Printf("Unsupported database engine: %s, falling back to SQLite", db.Engine)
		}
		cwd, err := os.Getwd()
		if err != nil {
			if cfg.IsDebug() {
				log.Printf("Error getting current working directory: %v", err)
			}
			cwd = "."
		}
		dbPath := filepath.Join(cwd, "db.sqlite3")
		dbURL = fmt.Sprintf("sqlite3://%s", dbPath)
	}
	if cfg.IsDebug() {
		log.Printf("Database URL for dbmate: %s", dbURL)
	}
	return dbURL
}

func EnsureDatabaseExists(cfg Config) error {
	if cfg.GetDatabaseConfig().Engine == "sqlite3" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting current working directory: %v", err)
		}
		dbPath := filepath.Join(cwd, cfg.GetDatabaseConfig().Name)
		return utils.EnsureFileExists(dbPath)
	}
	return nil
}

func InitConfig(cfg Config) {
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
	var settings SettingsStruct
	if err := viper.Unmarshal(&settings); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	// Override with environment variables if present
	if os.Getenv("DEBUG") != "" {
		cfg.SetDebug(os.Getenv("DEBUG") == "true")
	}
	if os.Getenv("DB_ENGINE") != "" {
		dbConfig := cfg.GetDatabaseConfig()
		dbConfig.Engine = os.Getenv("DB_ENGINE")
		cfg.SetDatabaseConfig(dbConfig)
	}
	if os.Getenv("DB_NAME") != "" {
		dbConfig := cfg.GetDatabaseConfig()
		dbConfig.Name = os.Getenv("DB_NAME")
		cfg.SetDatabaseConfig(dbConfig)
	}

	InitLogger(cfg.IsDebug())

	// Ensure the database file exists (for SQLite)
	if err := EnsureDatabaseExists(cfg); err != nil {
		log.Fatalf("Failed to ensure database exists: %v", err)
	}

	if cfg.IsDebug() {
		log.Println("Configuration initialized")
	}
}
